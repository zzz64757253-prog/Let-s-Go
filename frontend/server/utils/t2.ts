// server/utils/r2.ts
import { S3Client } from '@aws-sdk/client-s3'

export const useR2 = () => {
  const config = useRuntimeConfig()
  
  return new S3Client({
    region: 'auto',
    endpoint: config.r2Endpoint,
    credentials: {
      accessKeyId: config.r2AccessKeyId,
      secretAccessKey: config.r2SecretAccessKey,
    },
  })
}