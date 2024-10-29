variable "username" {
  type        = string
}

variable "admin_username"{
    type    = string
    default = "labuser"
}

variable "admin_password"{
    type    = string
    default = "S3cur3P4ssw0rd123!"
}

# VM Image details
variable "publisher" {
  type    = string
  default = "kali-linux"
}

variable "offer" {
  type    = string
  default = "kali"
}

variable "sku" {
  type    = string
  default = "kali-2024-2"
}

variable "vmversion" {
  type    = string
  default = "2024.2.0"
}

variable "kali_vm_size" {
  type    = string
  default = "Standard_D4ads_v5"
}

variable "vm_size" {
  type    = string
  default = "Standard_B2s"
}

# Ubuntu VM Image details
variable "ubupublisher" {
  type    = string
  default = "canonical"
}

variable "ubuoffer" {
  type    = string
  default = "0001-com-ubuntu-server-focal"
}

variable "ubusku" {
  type    = string
  default = "20_04-lts-gen2"
}

variable "ubuvmversion" {
  type    = string
  default = "20.04.202406140"
}