# Defining `email` contact point types

## Example usage

```yaml
apiVersion: k8s.kevingomez.fr/v1alpha1
kind: AlertManager
metadata:
  name: alertmanager-example
spec:
  contact_points:
    - name: Team A

      # Contact Team A by email
      contacts:
        - email: { to: ['team-a@unicorn.io', 'team-a-engineering@unicorn.io'] }

  routing:
    # ... omitted
```

## Reference

```yaml
apiVersion: k8s.kevingomez.fr/v1alpha1
kind: AlertManager
metadata:
  name: alertmanager-example
spec:
  # Alerts not matched by any of the routing rules will be sent to this contact point
  default_contact_point: 'Team A'

  # List of known contact points
  contact_points:
    - name: Team A
      contacts:
        - email:
            # List of recipients to reach out to.
            # Optional (but strongly suggested). Default: []
            to: ['team-a@unicorn.io', 'team-a-engineering@unicorn.io']

            # Send a single email to all recipients
            # Optional. Default: false
            single: false

            # Message to include with the email. You can use template variables.
            # Optional. Default: ''
            message: ''
```

## That was it!

[Return to the index to explore what you can do with DARK](../index.md)
