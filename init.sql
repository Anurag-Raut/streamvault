CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "User" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT UNIQUE NOT NULL,
    "profileImage" TEXT
);

CREATE TABLE "Password" (
    username TEXT PRIMARY KEY REFERENCES "User" (username),
    password TEXT NOT NULL
);

CREATE TABLE "Video" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    category TEXT ,
    thumbnail TEXT NOT NULL,
    "isStreaming" BOOLEAN NOT NULL,
    "userId" UUID NOT NULL REFERENCES "User" (id),
    "isVOD" BOOLEAN DEFAULT False,
    "isProcessed" BOOLEAN DEFAULT True,
    "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "views" INTEGER DEFAULT 0,
    "visibility" BOOLEAN DEFAULT True
);

CREATE TABLE "Like" (
    "videoId" UUID NOT NULL REFERENCES "Video" (id),
    "userId" UUID NOT NULL REFERENCES "User" (id),
    "isLike" BOOLEAN NOT NULL,
    PRIMARY KEY ("videoId", "userId")
);

CREATE TABLE "Subscription" (
    "creatorId" UUID NOT NULL REFERENCES "User" (id),
    "subscriberId" UUID NOT NULL REFERENCES "User" (id),
    "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("creatorId", "subscriberId")
);

CREATE TABLE "Comment" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "videoId" UUID NOT NULL REFERENCES "Video" (id),
    "userId" UUID NOT NULL REFERENCES "User" (id),
    text TEXT NOT NULL,
    "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);