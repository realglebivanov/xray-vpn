include $(TOPDIR)/rules.mk

PKG_NAME:=xray-vpn
PKG_VERSION:=0.1.0
PKG_RELEASE:=1

PKG_MAINTAINER:=Gleb Ivanov <realglebivanov@gmail.com>

include $(INCLUDE_DIR)/package.mk

define Package/xray-vpn
  SECTION:=net
  CATEGORY:=Network
  SUBMENU:=VPN
  TITLE:=Xray TUN VPN with automatic route management
  DEPENDS:=+ca-certificates +wget
  PKGARCH:=all
endef

define Package/xray-vpn/description
  Single daemon embedding xray-core with TUN mode.
  Fetches Russian IP ranges at startup, injects routing rules,
  manages fwmark-based split routing via netlink.
  Binaries are downloaded to /tmp on first start.
endef

XRAY_VPN_URL:=https://github.com/realglebivanov/xray-vpn/releases/download/v$(PKG_VERSION)

define Build/Compile
endef

define Package/xray-vpn/install
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_DIR) $(1)/etc/xray-vpn
	$(INSTALL_DIR) $(1)/usr/bin

	$(INSTALL_BIN) ./files/xray-vpn.init $(1)/etc/init.d/xray-vpn
	$(INSTALL_BIN) ./files/xray-vpn-download $(1)/usr/bin/xray-vpn-download
	$(INSTALL_DATA) /dev/null $(1)/etc/xray-vpn/state.json
endef

define Package/xray-vpn/conffiles
/etc/xray-vpn/state.json
endef

define Package/xray-vpn/postinst
#!/bin/sh
[ -n "$${IPKG_INSTROOT}" ] || {
	/etc/init.d/xray-vpn enable
	CRON_LINE="*/30 * * * * /tmp/xray-vpn refresh"
	grep -qF "$$CRON_LINE" /etc/crontabs/root 2>/dev/null || \
		echo "$$CRON_LINE" >> /etc/crontabs/root
	/etc/init.d/cron restart
}
endef

define Package/xray-vpn/prerm
#!/bin/sh
/etc/init.d/xray-vpn stop 2>/dev/null
/etc/init.d/xray-vpn disable
sed -i '\|/tmp/xray-vpn refresh|d' /etc/crontabs/root 2>/dev/null
/etc/init.d/cron restart 2>/dev/null
endef

$(eval $(call BuildPackage,xray-vpn))
