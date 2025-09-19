import { Module } from '@nestjs/common';
import { PrismaModule } from './prisma/prisma.module';
import { ServersModule } from './servers/servers.module';
import { ProvidersModule } from './providers/providers.module';
import { PeopleModule } from './people/people.module';
import { CostSnapshotsModule } from './cost-snapshots/cost-snapshots.module';

@Module({
  imports: [
    PrismaModule,
    ServersModule,
    ProvidersModule,
    PeopleModule,
    CostSnapshotsModule,
  ],
})
export class AppModule {}