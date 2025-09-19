import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
  ParseIntPipe,
  Query,
} from '@nestjs/common';
import { CostSnapshotsService } from './cost-snapshots.service';
import { CreateCostSnapshotDto } from './dto/create-cost-snapshot.dto';
import { UpdateCostSnapshotDto } from './dto/update-cost-snapshot.dto';

@Controller('cost-snapshots')
export class CostSnapshotsController {
  constructor(private readonly costSnapshotsService: CostSnapshotsService) {}

  @Post()
  create(@Body() createCostSnapshotDto: CreateCostSnapshotDto) {
    return this.costSnapshotsService.create(createCostSnapshotDto);
  }

  @Get()
  findAll(@Query('serverId', ParseIntPipe) serverId?: number) {
    if (serverId) {
      return this.costSnapshotsService.findByServer(serverId);
    }
    return this.costSnapshotsService.findAll();
  }

  @Get(':id')
  findOne(@Param('id', ParseIntPipe) id: number) {
    return this.costSnapshotsService.findOne(id);
  }

  @Patch(':id')
  update(
    @Param('id', ParseIntPipe) id: number,
    @Body() updateCostSnapshotDto: UpdateCostSnapshotDto,
  ) {
    return this.costSnapshotsService.update(id, updateCostSnapshotDto);
  }

  @Delete(':id')
  remove(@Param('id', ParseIntPipe) id: number) {
    return this.costSnapshotsService.remove(id);
  }
}