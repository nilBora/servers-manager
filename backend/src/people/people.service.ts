import { Injectable } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { CreatePersonDto } from './dto/create-person.dto';
import { UpdatePersonDto } from './dto/update-person.dto';

@Injectable()
export class PeopleService {
  constructor(private prisma: PrismaService) {}

  create(createPersonDto: CreatePersonDto) {
    return this.prisma.person.create({
      data: createPersonDto,
    });
  }

  findAll() {
    return this.prisma.person.findMany({
      orderBy: { name: 'asc' },
      include: {
        _count: {
          select: { serversOwned: true },
        },
      },
    });
  }

  findOne(id: number) {
    return this.prisma.person.findUnique({
      where: { id },
      include: {
        serversOwned: {
          orderBy: { name: 'asc' },
          include: {
            provider: true,
          },
        },
      },
    });
  }

  update(id: number, updatePersonDto: UpdatePersonDto) {
    return this.prisma.person.update({
      where: { id },
      data: updatePersonDto,
    });
  }

  remove(id: number) {
    return this.prisma.person.delete({
      where: { id },
    });
  }
}