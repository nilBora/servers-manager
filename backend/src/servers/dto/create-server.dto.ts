import { IsString, IsOptional, IsInt, IsEnum, IsNumber, IsDateString } from 'class-validator';
import { ServerStatus, ServerPurpose, BillingType } from '@prisma/client';

export class CreateServerDto {
  @IsString()
  name: string;

  @IsString()
  hostname: string;

  @IsOptional()
  @IsString()
  ipPublic?: string;

  @IsOptional()
  @IsString()
  ipPrivate?: string;

  @IsOptional()
  @IsInt()
  port?: number;

  @IsOptional()
  @IsString()
  username?: string;

  @IsOptional()
  @IsString()
  password?: string;

  @IsOptional()
  @IsString()
  sshKey?: string;

  @IsOptional()
  @IsEnum(ServerStatus)
  status?: ServerStatus;

  @IsOptional()
  @IsEnum(ServerPurpose)
  purpose?: ServerPurpose;

  @IsOptional()
  @IsEnum(BillingType)
  billingType?: BillingType;

  @IsOptional()
  @IsNumber()
  costMonthEstimated?: number;

  @IsOptional()
  @IsDateString()
  decommissionAt?: string;

  @IsOptional()
  @IsInt()
  providerId?: number;

  @IsOptional()
  @IsInt()
  ownerId?: number;

  @IsOptional()
  @IsString()
  os?: string;

  @IsOptional()
  @IsString()
  cpu?: string;

  @IsOptional()
  @IsString()
  ram?: string;

  @IsOptional()
  @IsString()
  storage?: string;

  @IsOptional()
  @IsString()
  location?: string;

  @IsOptional()
  @IsString()
  description?: string;

  @IsOptional()
  @IsString()
  tags?: string;

  @IsOptional()
  @IsString()
  account?: string;
}