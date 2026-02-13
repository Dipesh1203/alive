-- CreateEnum
CREATE TYPE "websiteStatus" AS ENUM ('Up', 'Down', 'Unknown');

-- CreateTable
CREATE TABLE "website" (
    "id" TEXT NOT NULL,
    "websiteName" TEXT NOT NULL,
    "url" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "website_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "region" (
    "regionId" TEXT NOT NULL,
    "regionName" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "region_pkey" PRIMARY KEY ("regionId")
);

-- CreateTable
CREATE TABLE "websiteTicks" (
    "id" SERIAL NOT NULL,
    "websiteId" TEXT NOT NULL,
    "upStatus" "websiteStatus" NOT NULL DEFAULT 'Unknown',
    "latency" INTEGER,
    "websiteRegionId" TEXT NOT NULL,

    CONSTRAINT "websiteTicks_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "websiteTicks" ADD CONSTRAINT "websiteTicks_websiteId_fkey" FOREIGN KEY ("websiteId") REFERENCES "website"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "websiteTicks" ADD CONSTRAINT "websiteTicks_websiteRegionId_fkey" FOREIGN KEY ("websiteRegionId") REFERENCES "region"("regionId") ON DELETE RESTRICT ON UPDATE CASCADE;
