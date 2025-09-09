---
title: "Task 2: Secure Cloud Infrastructure"
linkTitle: "Task 2: Secure Cloud Infrastructure"
weight: 20
---

|                            |    |  
|----------------------------| ----
| **Goal**                   | Learn how to lock down Access in Azure
| **Task**                   | Modify Azure NSG in terraform to only allow traffic from FortiAppSec Cloud
| **Verify task completion** | You should no longer be able to reach Juice Shop directly from your desktop.

This lab was originally bootstrapped with an ingress firewall rule which allows all ports and protocols from all sources (0.0.0.0/0).  This is not a best practice.  It is recommended, to only allow necessary ports and/or sources.  Now that we have onboarded our application, we want to ensure that the only device that can communicate with our application is FortiAppSec Cloud.

### Task 1: Modify Azure Network Firewall Rules

1. In the FortiAppSec Cloud UI, Copy the IPs which FortiAppSec Cloud will use to communicate with your application. From the FortiAppSec Cloud Applications page, select **Allow IP List** from the top of the page.  This will open a dialog showing Management and Scrubbing Center Addresses.  Click on **Copy to Clipboard**.  Paste these IPs into a text document and then click **Return**

  
   {{< figure src="allowIP-list.png" alt="allow-IP" >}}

2. In Azure cloud shell, verify you're in the terraform folder or navigate to it by typing ``` cd fortiweb-security-foundations-201/terraform/```

3. Make a copy of our ubuntu.tf file so that we can come back to it later if needed.
    - at the prompt, type ```cp ubuntu.tf ubuntu.tf.bak```

4. Use nano to open and edit the ubuntu terraform file

   {{< figure src="ubutf1.png" alt="ubutf1" >}}
   
   
   {{% notice info %}}
   In order to Navigate within nano, use the **up, down, left and right** arrow keys.  Use **backspace** to delete and type in the text you want to replace it with.  When you are ready to save type <kbd>ctrl</kbd>+<kbd>o</kbd> then  <kbd>enter</kbd> (to save to the same filename). Then type <kbd>ctrl</kbd>+<kbd>x</kbd> to exit.
   {{% /notice %}}

   {{< figure src="ubutf3.png" alt="ubutf3" >}}

5. Navigate to the **security rule** named **allow-juice-inbound**.  Note currently, the rule allows all source addresses.

   {{< figure src="ubutf2.png" alt="ubutf2" >}}

6. We are going to modify the **source_address_prefix** entry and replace it with the list of FortiWeb Cloud IPs captured in the step one above
   - For example: ```source_address_prefixes    = ["3.226.2.163", "3.123.68.65", "52.179.7.200", "20.127.74.103", "20.127.74.161", "20.127.74.143", "20.228.249.214", "52.179.3.225"]```
   - Save the file with <kbd>ctrl</kbd>+<kbd>o</kbd> then  <kbd>enter</kbd> (to save to the same filename)
   - Exit Nano with <kbd>ctrl</kbd>+<kbd>x</kbd> 
   - When you are done you can verify your changes with ```more ubuntu.tf```


      {{< figure src="ubutf4.png" alt="ubutf4" >}}

7. Now we will apply these changes by typing ```terraform apply -var="username=$(whoami)" --auto-approve```

8. When this is completed you will see terraform removed both security rules and added the new ones in their place.

   ```sh
   Terraform will perform the following actions:
   
     # azurerm_network_security_group.ubu-nsg will be updated in-place
     ~ resource "azurerm_network_security_group" "ubu-nsg" {
           id                  = "/subscriptions/02b50049-c444-416f-a126-3e4c815501ac/resourceGroups/web10-http101-workshop/providers/Microsoft.Network/networkSecurityGroups/web10-ubu_nsg"
           name                = "web10-ubu_nsg"
         ~ security_rule       = [
             - {
                 - access                                     = "Allow"
                 - destination_address_prefix                 = "*"
                 - destination_address_prefixes               = []
                 - destination_application_security_group_ids = []
                 - destination_port_range                     = "22"
                 - destination_port_ranges                    = []
                 - direction                                  = "Inbound"
                 - name                                       = "allow-ssh-inbound"
                 - priority                                   = 101
                 - protocol                                   = "Tcp"
                 - source_address_prefix                      = "*"
                 - source_address_prefixes                    = []
                 - source_application_security_group_ids      = []
                 - source_port_range                          = "*"
                 - source_port_ranges                         = []
                   # (1 unchanged attribute hidden)
               },
             - {
                 - access                                     = "Allow"
                 - destination_address_prefix                 = "*"
                 - destination_address_prefixes               = []
                 - destination_application_security_group_ids = []
                 - destination_port_range                     = "3000"
                 - destination_port_ranges                    = []
                 - direction                                  = "Inbound"
                 - name                                       = "allow-juice-inbound"
                 - priority                                   = 102
                 - protocol                                   = "Tcp"
                 - source_address_prefix                      = "*"
                 - source_address_prefixes                    = []
                 - source_application_security_group_ids      = []
                 - source_port_range                          = "*"
                 - source_port_ranges                         = []
                   # (1 unchanged attribute hidden)
               },
             + {
                 + access                                     = "Allow"
                 + destination_address_prefix                 = "*"
                 + destination_address_prefixes               = []
                 + destination_application_security_group_ids = []
                 + destination_port_range                     = "3000"
                 + destination_port_ranges                    = []
                 + direction                                  = "Inbound"
                 + name                                       = "allow-juice-inbound"
                 + priority                                   = 102
                 + protocol                                   = "Tcp"
                 + source_address_prefixes                    = [
                     + "3.123.68.65",
                     + "3.226.2.163",
                     + "34.138.149.79",
                     + "34.148.6.49",
                     + "34.74.199.185",
                     + "35.185.18.199",
                     + "35.227.112.86",
                     + "35.227.32.42",
                   ]
                 + source_application_security_group_ids      = []
                 + source_port_range                          = "*"
                 + source_port_ranges                         = []
                   # (2 unchanged attributes hidden)
               },
             + {
                 + access                                     = "Allow"
                 + destination_address_prefix                 = "*"
                 + destination_address_prefixes               = []
                 + destination_application_security_group_ids = []
                 + destination_port_range                     = "22"
                 + destination_port_ranges                    = []
                 + direction                                  = "Inbound"
                 + name                                       = "allow-ssh-inbound"
                 + priority                                   = 101
                 + protocol                                   = "Tcp"
                 + source_address_prefix                      = "*"
                 + source_address_prefixes                    = []
                 + source_application_security_group_ids      = []
                 + source_port_range                          = "*"
                 + source_port_ranges                         = []
               },
           ]
           tags                = {}
           # (2 unchanged attributes hidden)
       }
   
   Plan: 0 to add, 1 to change, 0 to destroy.
   azurerm_network_security_group.ubu-nsg: Modifying... [id=/subscriptions/02b50049-c444-416f-a126-3e4c815501ac/resourceGroups/web10-http101-workshop/providers/Microsoft.Network/networkSecurityGroups/web10-ubu_nsg]
   azurerm_network_security_group.ubu-nsg: Modifications complete after 2s [id=/subscriptions/02b50049-c444-416f-a126-3e4c815501ac/resourceGroups/web10-http101-workshop/providers/Microsoft.Network/networkSecurityGroups/web10-ubu_nsg]
   
   Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
   ```

9. You can verify this change in the Azure Portal as well.  
   - From the home screen in Azure, search in the top middle bar for  ```ubu_nsg```
   - You will find a **Network Security Group** with a name corresponding to your Azure Account ID like **web10-ubu_nsg**.  Click to view it
   
     {{< figure src="nsg1.png" alt="nsg1" >}}

10. You should be able to see the updated security rule.

   {{< figure src="nsg2.png" alt="nsg2" >}}

11. Now try to navigate to the Juice Shop Application from your laptop by typing ```http://<ubuntu ip>:3000``` in your favorite browser.
    - You should **NOT** be able to access Juice Shop Directly.
