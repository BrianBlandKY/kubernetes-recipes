On power outage we need to reset pihole.
- Perform the normal delete\create
- Then connnect to pihole terminal and set password

kubectl exec -ti earthshaker-pihole-f576b446f-pbkcl -c earthshaker-pihole -n shawshank bash
pihole -a -p
Enter New Password (Blank for no password):
  [✓] Password Removed