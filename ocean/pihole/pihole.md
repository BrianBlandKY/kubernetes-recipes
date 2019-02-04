On power outage we need to reset pihole.
- Perform the normal delete\create
- Then connnect to pihole terminal and set password

kubectl exec -ti POD -c pihole -n ocean bash
pihole -a -p
Enter New Password (Blank for no password):
  [✓] Password Removed