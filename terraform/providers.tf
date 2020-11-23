terraform {
  required_version = ">=0.12"
  backend "s3" {
    bucket = "myfiles.kazuya"
    key    = "terraform/power-phrase2.tfstate"
    region = "ap-northeast-1"
  }
}

provider aws {
  region = "ap-northeast-1"
}
