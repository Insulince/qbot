#!/bin/sh
set -e

# Ensure database directory exists
mkdir -p /var/lib/litefs

# Initialize database schema (idempotent)
sqlite3 /var/lib/litefs/qbot.db < /app/schema.sql

# Start LiteFS + bot
exec litefs mount -- /app/bot
