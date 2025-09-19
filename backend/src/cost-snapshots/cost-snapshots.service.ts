import { Injectable } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { CreateCostSnapshotDto } from './dto/create-cost-snapshot.dto';
import { UpdateCostSnapshotDto } from './dto/update-cost-snapshot.dto';

@Injectable()
export class CostSnapshotsService {
  constructor(private prisma: PrismaService) {}

  create(createCostSnapshotDto: CreateCostSnapshotDto) {
    const { month, ...data } = createCostSnapshotDto;
    return this.prisma.costSnapshot.create({
      data: {
        ...data,
        month: new Date(month),
      },
      include: {
        server: true,
      },
    });
  }

  findAll() {
    return this.prisma.costSnapshot.findMany({
      orderBy: { month: 'desc' },
      include: {
        server: {
          select: {
            id: true,
            name: true,
            hostname: true,
          },
        },
      },
    });
  }

  findByServer(serverId: number) {
    return this.prisma.costSnapshot.findMany({
      where: { serverId },
      orderBy: { month: 'desc' },
      include: {
        server: {
          select: {
            id: true,
            name: true,
            hostname: true,
          },
        },
      },
    });
  }

  findOne(id: number) {
    return this.prisma.costSnapshot.findUnique({
      where: { id },
      include: {
        server: true,
      },
    });
  }

  update(id: number, updateCostSnapshotDto: UpdateCostSnapshotDto) {
    const { month, ...data } = updateCostSnapshotDto;
    return this.prisma.costSnapshot.update({
      where: { id },
      data: {
        ...data,
        month: month ? new Date(month) : undefined,
      },
      include: {
        server: true,
      },
    });
  }

  remove(id: number) {
    return this.prisma.costSnapshot.delete({
      where: { id },
    });
  }
}