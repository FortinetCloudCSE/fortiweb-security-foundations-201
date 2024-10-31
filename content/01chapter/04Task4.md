---
title: "Task 4: Check Juice Shop"
linkTitle: "Task 4: Check Juice Shop"
weight: 40
---

|                            |    |  
|----------------------------| ----
| **Goal**                   | Verify that Juice Shop is working
| **Task**                   | Navigate to the public IP associated with Ubuntu
| **Verify task completion** | You should see the Juice Shop Home Page



### Start Kali RDP

From your Terraform Outputs in Task 2 you should have seen gotten the Public IP address of Ubnutu.  

- By default, Juice shop listens on port 3000.  In your favorite browser, type ```http://<ubuntu-ip>:3000``` 
- You should see a screen like below:

{{% notice warning %}} Depending on your browser, you will likely need to accept the self-signed certificate warnings.  {{% /notice %}}

![Juice Shop](js-initial.png)

- You can now proceed to the next module
