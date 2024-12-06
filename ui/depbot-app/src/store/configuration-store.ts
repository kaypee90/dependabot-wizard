import { makeAutoObservable } from "mobx";

class ConfigurationStore {
  yamlData: string = "";

  ecosystem: string = "";
  directory: string = "";
  schedule: string = "";
  reviewer: string = "";
  openPrLimit: string = "";
  label: string = "";

  showGithubInfo: boolean = false;

  constructor() {
    makeAutoObservable(this);
  }

  setPackageConfiguration(ecosystem: string, directory: string, schedule: string, reviewer: string, openPrLimit: string, label: string) {
    this.ecosystem = ecosystem;
    this.directory = directory;
    this.schedule = schedule;
    this.reviewer = reviewer;
    this.openPrLimit = openPrLimit;
    this.label = label;
  }

  saveYamlData(newYamlData: string) {
    if (newYamlData.trim()) {
      this.yamlData = newYamlData.trim() + "\n";
    }
  }

  setShowGithubInfo(showInfo: boolean) {
    this.showGithubInfo = showInfo;
  }
}

export const configStore = new ConfigurationStore();

