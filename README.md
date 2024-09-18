# GPG Terraform provider

The GPG provider allows to generate GPG encrypted message in ASCII-armored format using Terraform. It is a fork of https://github.com/invidian/terraform-provider-gpg.


## Configuring your repository

You can change settings in the repository by manupilating configuration file
in `.pallet/gitconfig.yaml`

The full definition of the api can be found [here][gitconfig-api-ref]

CI of the changes in the repositoryconfig is not configured yet. For now you
can use the [argocd app status][argocd-app-ref] to find the status and
potential errors.


[gitconfig-api-ref]: https://github.com/coopnorge/cloud-platform-apis/blob/main/cloud-platform-apis/templates/repositoryconfig.github.coop.no/definition.yaml
[argocd-app-ref]:  https://argocd.internal.coop/applications?search=pallet-terraform-provider-opengpg&showFavorites=false&proj=&sync=&autoSync=&health=&namespace=&cluster=&labels=
