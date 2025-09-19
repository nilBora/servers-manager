import { IsString, IsOptional } from 'class-validator';

export class CreateProviderDto {
  @IsString()
  name: string;

  @IsOptional()
  @IsString()
  consoleUrl?: string;

  @IsOptional()
  @IsString()
  notes?: string;
}