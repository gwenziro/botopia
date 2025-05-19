/**
 * Configuration Application
 *
 * Mengelola data dan logika untuk halaman konfigurasi Botopia
 */

document.addEventListener("alpine:init", () => {
  Alpine.data("configApp", () => ({
    isConnected: false,
    isGoogleConnected: false,
    connectedPhone: "",
    spreadsheetUrl: "",
    driveFolderUrl: "",
    name: "",
    isSaving: false,
    stats: {
      uptime: 0,
    },
    config: {
      commandPrefix: "!",
      logLevel: "INFO",
      webPort: 8080,
      webHost: "0.0.0.0",
      webAuthEnabled: false,
      webAuthUsername: "admin",
      webAuthPassword: "",
      spreadsheetId: "",
      driveFolderId: "",
      credentialsFile: "./service-account.json",
    },
    originalConfig: {},

    initialize() {
      console.log("Initializing config app");

      // Ambil data statistik dan status koneksi
      this.fetchStats();

      // Ambil konfigurasi saat ini
      this.fetchConfig();

      // Ambil status konfigurasi yang lebih detail
      this.fetchConfigStatus();

      // Tambahkan event listener untuk konfirmasi sebelum pengguna meninggalkan halaman
      window.addEventListener("beforeunload", (e) => {
        if (this.hasChanges()) {
          e.preventDefault();
          e.returnValue =
            "Anda memiliki perubahan yang belum disimpan. Yakin ingin meninggalkan halaman?";
        }
      });
    },

    fetchStats() {
      fetch("/api/stats")
        .then((response) => response.json())
        .then((data) => {
          this.isConnected = data.isConnected;
          this.connectedPhone = data.phone || "";
          this.stats = {
            uptime: data.uptime || 0,
          };

          if (data.name && data.name !== "WhatsApp User") {
            this.name = data.name;
          } else if (data.pushName && data.pushName !== "") {
            this.name = data.pushName;
          } else {
            this.name = data.phone || "WhatsApp User";
          }
        })
        .catch((error) => {
          console.error("Error fetching stats:", error);
          this.showErrorToast("Gagal memuat statistik sistem");
        });
    },

    fetchConfig() {
      fetch("/api/config")
        .then((response) => response.json())
        .then((data) => {
          console.log("Received config data:", data);

          // Update local config with server data
          this.config = {
            ...this.config,
            commandPrefix: data.commandPrefix || "!",
            logLevel: data.logLevel || "INFO",
            webPort: data.webPort || 8080,
            webHost: data.webHost || "0.0.0.0",
            webAuthEnabled: data.webAuthEnabled || false,
            webAuthUsername: data.webAuthUsername || "admin",
            spreadsheetId: data.googleSheets?.spreadsheetID || "",
            driveFolderId: data.googleSheets?.driveFolderID || "",
            credentialsFile:
              data.googleSheets?.credentialsFile || "./service-account.json",
          };

          // Store original config for detecting changes
          this.originalConfig = JSON.parse(JSON.stringify(this.config));
        })
        .catch((error) => {
          console.error("Error fetching config:", error);
          this.showErrorToast("Gagal memuat konfigurasi sistem");
        });
    },

    fetchConfigStatus() {
      fetch("/api/config/status")
        .then((response) => response.json())
        .then((data) => {
          console.log("Received config status:", data);

          // Set Google connection status
          this.isGoogleConnected = data.googleApi?.configured || false;

          // Set spreadsheet URL if available
          if (data.spreadsheetUrl) {
            this.spreadsheetUrl = data.spreadsheetUrl;
          }

          // Set drive folder URL if available
          if (data.driveFolderUrl) {
            this.driveFolderUrl = data.driveFolderUrl;
          }
        })
        .catch((error) => {
          console.error("Error fetching config status:", error);
        });
    },

    setCommandPrefix(prefix) {
      this.config.commandPrefix = prefix;
    },

    saveConfig() {
      if (this.isSaving) return;

      this.isSaving = true;
      const saveBtn = document.querySelector(".c-save-button");
      if (saveBtn) {
        const originalHtml = saveBtn.innerHTML;
        saveBtn.innerHTML =
          '<i class="fas fa-spinner fa-spin"></i> Menyimpan...';
        saveBtn.disabled = true;

        // Prepare the data to send to server
        const configData = {
          commandPrefix: this.config.commandPrefix,
          logLevel: this.config.logLevel,
          webPort: this.config.webPort,
          webHost: this.config.webHost,
          webAuthEnabled: this.config.webAuthEnabled,
          webAuthUsername: this.config.webAuthUsername,
          googleSheets: {
            spreadsheetID: this.config.spreadsheetId,
            driveFolderID: this.config.driveFolderId,
            credentialsFile: this.config.credentialsFile,
          },
        };

        // If password was modified, include it
        if (
          this.config.webAuthPassword &&
          this.config.webAuthPassword !== "********"
        ) {
          configData.webAuthPassword = this.config.webAuthPassword;
        }

        fetch("/api/config", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(configData),
        })
          .then((response) => {
            if (!response.ok) {
              throw new Error("Network response was not ok");
            }
            return response.json();
          })
          .then((data) => {
            if (data.success) {
              this.showSuccessToast("Konfigurasi berhasil disimpan");

              // Update original config to reflect saved state
              this.originalConfig = JSON.parse(JSON.stringify(this.config));

              // Refresh configuration status
              this.fetchConfigStatus();
            } else {
              throw new Error(data.message || "Unknown error");
            }
          })
          .catch((error) => {
            console.error("Error saving config:", error);
            this.showErrorToast(
              `Gagal menyimpan konfigurasi: ${error.message}`
            );
          })
          .finally(() => {
            this.isSaving = false;
            if (saveBtn) {
              saveBtn.innerHTML = originalHtml;
              saveBtn.disabled = false;
            }
          });
      }
    },

    hasChanges() {
      return (
        JSON.stringify(this.originalConfig) !== JSON.stringify(this.config)
      );
    },

    formatUptime(seconds) {
      if (!seconds || seconds <= 0) return "Tidak aktif";

      // Konversi detik ke format yang lebih ramah
      const days = Math.floor(seconds / 86400);
      const hours = Math.floor((seconds % 86400) / 3600);
      const minutes = Math.floor((seconds % 3600) / 60);

      let result = [];
      if (days > 0) {
        result.push(`${days} hari`);
      }
      if (hours > 0 || days > 0) {
        result.push(`${hours} jam`);
      }
      if (minutes > 0 || (hours > 0 && days === 0) || days > 0) {
        result.push(`${minutes} menit`);
      } else {
        result.push(`${seconds % 60} detik`);
      }

      return result.join(" ");
    },

    showSuccessToast(message) {
      if (typeof showToast === "function") {
        showToast("success", message);
      } else {
        alert(message);
      }
    },

    showErrorToast(message) {
      if (typeof showToast === "function") {
        showToast("error", message);
      } else {
        alert("Error: " + message);
      }
    },
  }));
});
