<?php

namespace SatBridge;

use PhpCfdi\CfdiSatScraper\SatScraper;
use PhpCfdi\CfdiSatScraper\Sessions\Ciec\CiecSessionManager;
use PhpCfdi\CfdiSatScraper\Sessions\Fiel\FielSessionManager;
use PhpCfdi\ImageCaptchaResolver\Resolvers\AntiCaptchaResolver;
use PhpCfdi\CfdiSatScraper\QueryByFilters\QueryByFiltersTranslator;
use PhpCfdi\CfdiSatScraper\Filters\DownloadType;
use PhpCfdi\CfdiSatScraper\Filters\Options\RfcOption;
use PhpCfdi\CfdiSatScraper\Filters\Options\StatesVoucherOption;
use PhpCfdi\CfdiSatScraper\Filters\Options\ComplementsOption;
use Psr\Log\LoggerInterface;
use DateTimeImmutable;
use Exception;

class SatScraperService
{
    private LoggerInterface $logger;
    private string $storagePath;

    public function __construct(LoggerInterface $logger)
    {
        $this->logger = $logger;
        $this->storagePath = $_ENV['STORAGE_PATH'] ?? './storage/cfdis';
        
        if (!is_dir($this->storagePath)) {
            mkdir($this->storagePath, 0755, true);
        }
    }

    /**
     * Download CFDIs based on filters
     */
    public function downloadCFDIs(array $params): array
    {
        try {
            $this->logger->info('Starting CFDI download', ['params' => $params]);

            // Validate required parameters
            $this->validateDownloadParams($params);

            // Create SAT scraper instance
            $scraper = $this->createScraper($params);

            // Build query filters
            $downloadType = $this->getDownloadType($params['tipo_cfdi'] ?? 'emitidos');
            $dateStart = new DateTimeImmutable($params['fecha_inicio']);
            $dateEnd = new DateTimeImmutable($params['fecha_fin']);

            // Create query
            $query = QueryByFiltersTranslator::createQueryByFilters(
                $dateStart,
                $dateEnd,
                $downloadType
            );

            // Apply additional filters
            if (isset($params['rfc_emisor'])) {
                $query = $query->setRfcEmisor(new RfcOption($params['rfc_emisor']));
            }
            if (isset($params['rfc_receptor'])) {
                $query = $query->setRfcReceptor(new RfcOption($params['rfc_receptor']));
            }
            if (isset($params['estado'])) {
                $query = $query->setStateVoucher(StatesVoucherOption::create($params['estado']));
            }
            if (isset($params['complemento'])) {
                $query = $query->setComplement(ComplementsOption::create($params['complemento']));
            }

            // Execute query
            $this->logger->info('Executing SAT query');
            $list = $scraper->listByPeriod($query);

            $downloadedCfdis = [];
            $errors = [];

            // Download each CFDI
            foreach ($list as $metadata) {
                try {
                    $uuid = $metadata->uuid();
                    $this->logger->info('Downloading CFDI', ['uuid' => $uuid]);

                    // Download XML
                    $xml = $scraper->downloadXml($metadata);
                    
                    // Save XML to file
                    $fileName = $this->saveXmlToFile($uuid, $xml);

                    $downloadedCfdis[] = [
                        'uuid' => $uuid,
                        'rfc_emisor' => $metadata->rfcEmisor(),
                        'rfc_receptor' => $metadata->rfcReceptor(),
                        'fecha_emision' => $metadata->fechaEmision()->format('Y-m-d'),
                        'monto_total' => $metadata->total(),
                        'status_sat' => $metadata->estadoComprobante(),
                        'tipo_cfdi' => $params['tipo_cfdi'] ?? 'emitidos',
                        'archivo_xml' => $fileName,
                        'xml_content' => base64_encode($xml)
                    ];

                    // Check limit
                    if (count($downloadedCfdis) >= ($_ENV['MAX_QUERIES_PER_REQUEST'] ?? 500)) {
                        $this->logger->warning('Reached maximum downloads per request');
                        break;
                    }
                } catch (Exception $e) {
                    $errors[] = [
                        'uuid' => $metadata->uuid(),
                        'error' => $e->getMessage()
                    ];
                    $this->logger->error('Failed to download CFDI', [
                        'uuid' => $metadata->uuid(),
                        'error' => $e->getMessage()
                    ]);
                }
            }

            return [
                'success' => true,
                'total_found' => count($list),
                'total_downloaded' => count($downloadedCfdis),
                'cfdis' => $downloadedCfdis,
                'errors' => $errors
            ];

        } catch (Exception $e) {
            $this->logger->error('CFDI download failed', ['error' => $e->getMessage()]);
            return [
                'success' => false,
                'error' => $e->getMessage()
            ];
        }
    }

    /**
     * Query metadata only without downloading
     */
    public function queryMetadata(array $params): array
    {
        try {
            $this->logger->info('Querying CFDI metadata', ['params' => $params]);

            $scraper = $this->createScraper($params);
            $downloadType = $this->getDownloadType($params['tipo_cfdi'] ?? 'emitidos');
            $dateStart = new DateTimeImmutable($params['fecha_inicio']);
            $dateEnd = new DateTimeImmutable($params['fecha_fin']);

            $query = QueryByFiltersTranslator::createQueryByFilters(
                $dateStart,
                $dateEnd,
                $downloadType
            );

            $list = $scraper->listByPeriod($query);

            $metadata = [];
            foreach ($list as $item) {
                $metadata[] = [
                    'uuid' => $item->uuid(),
                    'rfc_emisor' => $item->rfcEmisor(),
                    'rfc_receptor' => $item->rfcReceptor(),
                    'fecha_emision' => $item->fechaEmision()->format('Y-m-d'),
                    'monto_total' => $item->total(),
                    'status_sat' => $item->estadoComprobante(),
                ];
            }

            return [
                'success' => true,
                'total_found' => count($metadata),
                'metadata' => $metadata
            ];

        } catch (Exception $e) {
            $this->logger->error('Metadata query failed', ['error' => $e->getMessage()]);
            return [
                'success' => false,
                'error' => $e->getMessage()
            ];
        }
    }

    /**
     * Download specific CFDI by UUID
     */
    public function downloadByUUID(array $params): array
    {
        try {
            $uuid = $params['uuid'];
            $this->logger->info('Downloading CFDI by UUID', ['uuid' => $uuid]);

            $scraper = $this->createScraper($params);

            // Note: PhpCfdi doesn't have direct UUID download, 
            // so we need to search for it in a date range
            if (!isset($params['fecha_inicio']) || !isset($params['fecha_fin'])) {
                throw new Exception('Date range required for UUID search');
            }

            $downloadType = $this->getDownloadType($params['tipo_cfdi'] ?? 'emitidos');
            $dateStart = new DateTimeImmutable($params['fecha_inicio']);
            $dateEnd = new DateTimeImmutable($params['fecha_fin']);

            $query = QueryByFiltersTranslator::createQueryByFilters(
                $dateStart,
                $dateEnd,
                $downloadType
            );

            $list = $scraper->listByPeriod($query);

            foreach ($list as $metadata) {
                if ($metadata->uuid() === $uuid) {
                    $xml = $scraper->downloadXml($metadata);
                    $fileName = $this->saveXmlToFile($uuid, $xml);

                    return [
                        'success' => true,
                        'cfdi' => [
                            'uuid' => $uuid,
                            'rfc_emisor' => $metadata->rfcEmisor(),
                            'rfc_receptor' => $metadata->rfcReceptor(),
                            'fecha_emision' => $metadata->fechaEmision()->format('Y-m-d'),
                            'monto_total' => $metadata->total(),
                            'status_sat' => $metadata->estadoComprobante(),
                            'archivo_xml' => $fileName,
                            'xml_content' => base64_encode($xml)
                        ]
                    ];
                }
            }

            return [
                'success' => false,
                'error' => 'CFDI not found with UUID: ' . $uuid
            ];

        } catch (Exception $e) {
            $this->logger->error('UUID download failed', ['error' => $e->getMessage()]);
            return [
                'success' => false,
                'error' => $e->getMessage()
            ];
        }
    }

    private function createScraper(array $params): SatScraper
    {
        $authType = $params['auth_type'] ?? 'ciec';

        if ($authType === 'fiel') {
            return $this->createFielScraper($params);
        } else {
            return $this->createCiecScraper($params);
        }
    }

    private function createCiecScraper(array $params): SatScraper
    {
        if (!isset($params['rfc']) || !isset($params['ciec'])) {
            throw new Exception('RFC and CIEC are required for CIEC authentication');
        }

        $captchaResolver = $this->createCaptchaResolver();
        $sessionManager = new CiecSessionManager($params['rfc'], $params['ciec'], $captchaResolver);

        return new SatScraper($sessionManager);
    }

    private function createFielScraper(array $params): SatScraper
    {
        if (!isset($params['rfc']) || !isset($params['certificado_cer']) || 
            !isset($params['clave_privada_key']) || !isset($params['password_key'])) {
            throw new Exception('RFC, certificate, private key, and password are required for FIEL authentication');
        }

        // Decode base64 files
        $cerContents = base64_decode($params['certificado_cer']);
        $keyContents = base64_decode($params['clave_privada_key']);

        // Save temporarily
        $cerPath = $this->storagePath . '/temp_' . uniqid() . '.cer';
        $keyPath = $this->storagePath . '/temp_' . uniqid() . '.key';
        
        file_put_contents($cerPath, $cerContents);
        file_put_contents($keyPath, $keyContents);

        try {
            $captchaResolver = $this->createCaptchaResolver();
            $sessionManager = new FielSessionManager(
                $cerPath,
                $keyPath,
                $params['password_key'],
                $captchaResolver
            );

            $scraper = new SatScraper($sessionManager);

            // Clean up temp files
            unlink($cerPath);
            unlink($keyPath);

            return $scraper;
        } catch (Exception $e) {
            // Clean up temp files on error
            if (file_exists($cerPath)) unlink($cerPath);
            if (file_exists($keyPath)) unlink($keyPath);
            throw $e;
        }
    }

    private function createCaptchaResolver()
    {
        $resolverToken = $_ENV['CAPTCHA_RESOLVER_TOKEN'] ?? '';
        
        if (empty($resolverToken)) {
            throw new Exception('Captcha resolver token not configured');
        }

        // Use AntiCaptcha or BoxFactura resolver
        return new AntiCaptchaResolver($resolverToken);
    }

    private function getDownloadType(string $tipo): DownloadType
    {
        return match($tipo) {
            'emitidos', 'ingresos' => DownloadType::issued(),
            'recibidos', 'egresos' => DownloadType::received(),
            default => DownloadType::issued()
        };
    }

    private function saveXmlToFile(string $uuid, string $xmlContent): string
    {
        $fileName = $uuid . '.xml';
        $filePath = $this->storagePath . '/' . $fileName;
        
        file_put_contents($filePath, $xmlContent);
        
        return $fileName;
    }

    private function validateDownloadParams(array $params): void
    {
        $required = ['auth_type', 'rfc', 'fecha_inicio', 'fecha_fin'];
        
        foreach ($required as $field) {
            if (!isset($params[$field]) || empty($params[$field])) {
                throw new Exception("Missing required parameter: $field");
            }
        }

        if ($params['auth_type'] === 'ciec' && !isset($params['ciec'])) {
            throw new Exception('CIEC is required when auth_type is ciec');
        }

        if ($params['auth_type'] === 'fiel' && 
            (!isset($params['certificado_cer']) || !isset($params['clave_privada_key']) || !isset($params['password_key']))) {
            throw new Exception('Certificate, private key, and password are required when auth_type is fiel');
        }
    }
}

