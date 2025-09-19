import { Module } from '@nestjs/common';
import { CostSnapshotsService } from './cost-snapshots.service';
import { CostSnapshotsController } from './cost-snapshots.controller';

@Module({
  controllers: [CostSnapshotsController],
  providers: [CostSnapshotsService],
})
export class CostSnapshotsModule {}