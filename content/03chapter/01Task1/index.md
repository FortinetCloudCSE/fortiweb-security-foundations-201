---
title: "Task 1: Perform a simple SQL injection attack"
linkTitle: "Task 1: Perform a simple SQL injection attack"
weight: 10
---


|                            |    |  
|----------------------------| ----
| **Goal**                   | Perform SQL injection attack
| **Task**                   | Perform a simple browser based attack
| **Verify task completion** | You should see that the SQLi attack is accepted by Juice Shop.


1. Log into Kali linux: ```https://{{Kali IP}}:8443```

2. To avoid any DNS problems during this workshop, we'll create a static hosts file entry on the Kali Box to resolve our FortiAppSec Cloud protected application
   - Open the terminal emulator by clicking on the black box at the bottom of the Kali Home screen. At the prompt, type:

    ```sh
    bash
    sudo nano /etc/hosts
    ```

3. When the host file opens, add the following 2 lines to the bottom of the file, and save it.
    - Be sure to substitute your FortiAppSec Student ID in the fields   
    - To save the entries use: <kbd>ctrl</kbd>+<kbd>o</kbd> then  <kbd>enter</kbd> (to save to the same filename). 
    - To exit Nano: type <kbd>ctrl</kbd>+<kbd>x</kbd>

       ```
       <replace this value with IP address you have from FortiAppSec>    <FortiAppSecStudentID>.fwebtraincse.com
       <replace this value with IP address you have from FortiAppSec>   <FortiAppSecStudentID>.fwebtraincse.com
      ```    
      {{< figure src="hosts.png" alt="Hosts" >}}
 
4. Navigate to the Firefox browser (located at the top of Kali desktop) and enter our FortiAppSec Cloud Protected Juice Shop URL into the navigation bar ```https://<FortiAppSecStudentID>.fwebtraincse.com```.  Accept warnings and proceed to the application

    {{< figure src="studhome.png" alt="stud-home" >}}

5. Let’s perform a very simple SQLi attack. To perform a SQLi attack append ```?name=' OR 'x'='x``` to your URL.  Be sure that you use **YOUR NUMBER**.

   - For example (**be sure to use your studentId**)
     - ```https://669.fwebtraincse.com/?name=' OR 'x'='x```
    {{% notice info %}}
    The attack will go through, and you will see the Juice Shop Home page
    {{% /notice %}}

