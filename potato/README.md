# Install Potato Head Via Helm

A tiny multi-service app built specifically for chaos experiments. Very fast to deploy and safe to break.

## How to install

1. Clone the repo

```bash
git clone https://github.com/podtato-head/podtato-head.git
```

2. Open repo directory

```bash
cd podtato-head
```

3. Run helm install command

```bash
helm install \
--namespace=podtato \
--create-namespace \
podtato-head ./delivery/chart
```
