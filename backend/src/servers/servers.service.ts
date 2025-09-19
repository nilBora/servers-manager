import { Injectable } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { CreateServerDto } from './dto/create-server.dto';
import { UpdateServerDto } from './dto/update-server.dto';

@Injectable()
export class ServersService {
  constructor(private prisma: PrismaService) {}

  create(createServerDto: CreateServerDto) {
    const { decommissionAt, ...data } = createServerDto;
    return this.prisma.server.create({
      data: {
        ...data,
        decommissionAt: decommissionAt ? new Date(decommissionAt) : undefined,
      },
      include: {
        provider: true,
        owner: true,
        costSnapshots: true,
      },
    });
  }

  findAll() {
    return this.prisma.server.findMany({
      orderBy: { createdAt: 'desc' },
      include: {
        provider: true,
        owner: true,
        costSnapshots: {
          orderBy: { month: 'desc' },
          take: 1,
        },
      },
    });
  }

  findOne(id: number) {
    return this.prisma.server.findUnique({
      where: { id },
      include: {
        provider: true,
        owner: true,
        costSnapshots: {
          orderBy: { month: 'desc' },
        },
      },
    });
  }

  update(id: number, updateServerDto: UpdateServerDto) {
    const { decommissionAt, ...data } = updateServerDto;
    return this.prisma.server.update({
      where: { id },
      data: {
        ...data,
        decommissionAt: decommissionAt ? new Date(decommissionAt) : undefined,
      },
      include: {
        provider: true,
        owner: true,
        costSnapshots: true,
      },
    });
  }

  remove(id: number) {
    return this.prisma.server.delete({
      where: { id },
    });
  }
}