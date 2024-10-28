import { makeAutoObservable } from "mobx";

class ConfigurationStore {
  yamlData: string = "";
  constructor() {
    makeAutoObservable(this);
  }

  saveYamlData(newYamlData: string) {
    if (newYamlData.trim()) {
      this.yamlData = newYamlData.trim() + "\n";
    }
  }
}

export const configStore = new ConfigurationStore();

