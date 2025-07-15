#!/bin/sh

# Set your license key (env preferred)
LICENSE_KEY="${MAXMIND_LICENSE_KEY:-}"

# GeoIP edition to download
DB_EDITION="GeoLite2-City"
OUTPUT_DIR="geoip"
DB_FILE="${OUTPUT_DIR}/${DB_EDITION}.mmdb"
ARCHIVE_FILE="./${DB_EDITION}.tar.gz"
EXTRACT_DIR="./${DB_EDITION}_tmp"

# Check if license key is provided
if [ -z "$LICENSE_KEY" ]; then
  echo "ERROR: MAXMIND_LICENSE_KEY is not set."
  exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"
rm -rf "$EXTRACT_DIR"

echo "Downloading ${DB_EDITION} database..."
curl -s -L \
  "https://download.maxmind.com/app/geoip_download?edition_id=${DB_EDITION}&license_key=${LICENSE_KEY}&suffix=tar.gz" \
  -o "$ARCHIVE_FILE"

if [ $? -ne 0 ]; then
  echo "ERROR: Failed to download database archive"
  exit 1
fi

echo "Extracting database archive..."
mkdir "$EXTRACT_DIR"
tar -xzf "$ARCHIVE_FILE" -C "$EXTRACT_DIR"

# Find .mmdb file inside the extracted folder
MMDB_PATH=$(find "$EXTRACT_DIR" -name "*.mmdb" | head -n 1)

if [ -f "$MMDB_PATH" ]; then
  mv "$MMDB_PATH" "$DB_FILE"
  echo "Database updated: $DB_FILE"
else
  echo "ERROR: .mmdb file not found"
  rm -rf "$EXTRACT_DIR" "$ARCHIVE_FILE"
  exit 1
fi

# Clean up
rm -rf "$EXTRACT_DIR" "$ARCHIVE_FILE"

exit 0
