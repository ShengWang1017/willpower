#!/bin/bash
# Willpower Forge - Service Uninstallation Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="willpower-forge"
INSTALL_DIR="/opt/${APP_NAME}"
DATA_DIR="/var/lib/${APP_NAME}"
SERVICE_FILE="/etc/systemd/system/${APP_NAME}.service"

echo -e "${YELLOW}=====================================${NC}"
echo -e "${YELLOW}Willpower Forge Service Uninstaller${NC}"
echo -e "${YELLOW}=====================================${NC}"
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}Error: This script must be run as root (use sudo)${NC}"
   exit 1
fi

# Check if service exists
if [ ! -f "$SERVICE_FILE" ]; then
    echo -e "${YELLOW}Warning: Service file not found. The service may not be installed.${NC}"
    echo "Continue anyway? (y/n)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        exit 0
    fi
fi

echo -e "${RED}WARNING: This will remove the Willpower Forge service.${NC}"
echo ""
echo "The following will be removed:"
echo "  - Service: $SERVICE_FILE"
echo "  - Application: $INSTALL_DIR"
echo ""
echo -e "${YELLOW}The following will be KEPT (contains your data):${NC}"
echo "  - Data directory: $DATA_DIR"
echo "  - Database: $DATA_DIR/willpower.db"
echo ""
read -p "Do you want to continue? (y/n) " -r REPLY
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Uninstallation cancelled."
    exit 0
fi

# Stop and disable service
if systemctl is-active --quiet $APP_NAME; then
    echo -e "${GREEN}Stopping service...${NC}"
    systemctl stop $APP_NAME
fi

if systemctl is-enabled --quiet $APP_NAME 2>/dev/null; then
    echo -e "${GREEN}Disabling service...${NC}"
    systemctl disable $APP_NAME
fi

# Remove service file
if [ -f "$SERVICE_FILE" ]; then
    echo -e "${GREEN}Removing service file...${NC}"
    rm -f "$SERVICE_FILE"
    systemctl daemon-reload
fi

# Remove application directory
if [ -d "$INSTALL_DIR" ]; then
    echo -e "${GREEN}Removing application directory...${NC}"
    rm -rf "$INSTALL_DIR"
fi

# Ask about data directory
echo ""
echo -e "${YELLOW}Data directory: $DATA_DIR${NC}"
read -p "Do you want to remove the data directory (including database)? (y/n) " -r REPLY
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [ -d "$DATA_DIR" ]; then
        # Backup before deletion
        BACKUP_FILE="$HOME/willpower-forge-backup-$(date +%Y%m%d-%H%M%S).tar.gz"
        echo -e "${GREEN}Creating backup at: $BACKUP_FILE${NC}"
        tar -czf "$BACKUP_FILE" -C "$DATA_DIR" . 2>/dev/null || true

        echo -e "${GREEN}Removing data directory...${NC}"
        rm -rf "$DATA_DIR"

        echo -e "${GREEN}Backup created at: $BACKUP_FILE${NC}"
    fi
else
    echo -e "${GREEN}Data directory preserved at: $DATA_DIR${NC}"
fi

echo ""
echo -e "${GREEN}=====================================${NC}"
echo -e "${GREEN}Uninstallation completed!${NC}"
echo -e "${GREEN}=====================================${NC}"
echo ""
