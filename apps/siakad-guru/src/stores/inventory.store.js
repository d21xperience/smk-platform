import { defineStore } from 'pinia';
import { InventoryService } from '@/services/InventoryService';
import { InventoryItem } from '@/domain/inventory/models/InventoryItem';
import { AssetIssue } from '@/domain/inventory/models/AssetIssue';

export const useInventoryStore = defineStore('inventory', {
  state: () => ({
    assets: [], // Array of InventoryItem
    issues: [], // Array of AssetIssue (riwayat)
    currentBorrows: [], // Array of AssetIssue (sedang dipinjam)
    isLoading: false,
    error: null,
  }),
  getters: {
    availableAssets: (state) => state.assets.filter(a => a.isAvailable()),
    borrowedAssets: (state) => state.assets.filter(a => a.isBorrowed()),
    damagedAssets: (state) => state.assets.filter(a => a.isDamaged()),
    getAssetById: (state) => (id) => state.assets.find(a => a.id === id),
  },
  actions: {
    async fetchAssets(filters = {}) {
      this.isLoading = true;
      this.error = null;
      try {
        const data = await InventoryService.getAssets(filters);
        this.assets = data.map(a => new InventoryItem(a));
        return this.assets;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.isLoading = false;
      }
    },

    async fetchIssues(assetId = null) {
      this.isLoading = true;
      try {
        const data = await InventoryService.getIssues(assetId);
        this.issues = data.map(i => new AssetIssue(i));
        this.currentBorrows = this.issues.filter(i => i.isBorrowed());
        return this.issues;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.isLoading = false;
      }
    },

    async borrowAsset(assetId, payload) {
      this.isLoading = true;
      try {
        const data = await InventoryService.borrowAsset(assetId, payload);
        const newIssue = new AssetIssue(data);
        this.issues.push(newIssue);
        // Update asset status di list
        const asset = this.assets.find(a => a.id === assetId);
        if (asset) {
          asset.markAsBorrowed();
        }
        this.currentBorrows = this.issues.filter(i => i.isBorrowed());
        return newIssue;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.isLoading = false;
      }
    },

    async returnAsset(issueId) {
      this.isLoading = true;
      try {
        const data = await InventoryService.returnAsset(issueId);
        const updatedIssue = new AssetIssue(data);
        // Update di issues
        const index = this.issues.findIndex(i => i.id === issueId);
        if (index !== -1) this.issues[index] = updatedIssue;
        // Update asset status
        const asset = this.assets.find(a => a.id === updatedIssue.assetId);
        if (asset) {
          asset.markAsReturned();
        }
        this.currentBorrows = this.issues.filter(i => i.isBorrowed());
        return updatedIssue;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.isLoading = false;
      }
    },

    async reportDamage(assetId, description) {
      this.isLoading = true;
      try {
        const data = await InventoryService.reportDamage(assetId, description);
        // Update local asset
        const asset = this.assets.find(a => a.id === assetId);
        if (asset) {
          asset.markAsDamaged(description);
        }
        return data;
      } catch (err) {
        this.error = err.message;
        throw err;
      } finally {
        this.isLoading = false;
      }
    },

    clear() {
      this.assets = [];
      this.issues = [];
      this.currentBorrows = [];
      this.error = null;
    },
  },
});
