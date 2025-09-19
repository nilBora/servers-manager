import { PartialType } from '@nestjs/mapped-types';
import { CreateCostSnapshotDto } from './create-cost-snapshot.dto';

export class UpdateCostSnapshotDto extends PartialType(CreateCostSnapshotDto) {}