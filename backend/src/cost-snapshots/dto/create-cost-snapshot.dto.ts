import { IsInt, IsDateString, IsNumber, IsOptional, IsString } from 'class-validator';

export class CreateCostSnapshotDto {
  @IsInt()
  serverId: number;

  @IsDateString()
  month: string;

  @IsNumber()
  costMonth: number;

  @IsOptional()
  @IsString()
  source?: string;
}