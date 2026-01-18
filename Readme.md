# Technologies

- Golang
- Terraform
- Localstack
- AWS S3
- AWS Glue
- AWS Athena

Todo:

- [ ] Create cronjob to convert Json into Parquet
- [ ] Change Glue + Athena to fetch new Parquet data
- [ ] Setup Terraform for local development
- [ ] Deploy with K8s

Run GitHub pipelines locally with Act:

```bash
act -W .github/workflows/ci.yml \
-P ubuntu-latest=catthehacker/ubuntu:act-latest \
--secret-file .env
```
