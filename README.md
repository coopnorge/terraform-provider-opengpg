# GPG Terraform provider

The GPG provider allows to generate GPG encrypted message in ASCII-armored format
using Terraform. It is a fork of [invidian/terraform-provider-gpg](https://github.com/invidian/terraform-provider-gpg).
The provider, and it's resources has been renamed from `gpg` to `opengpg`, to make
a clearer distinction between the two providers.

The provider has been forked to add support for newer encryption-algorithms,
such as [Curve25519](https://en.wikipedia.org/wiki/Curve25519).

The original provider was using `golang.org/x/crypto/openpgp`, but that has been
marked frozen and deprecated, and it does not support Curve25519.
The Go team is promoting using community-forks such as [keybase/go-crypto](https://github.com/keybase/go-crypto)
or [ProtonMail/go-crypto](https://github.com/ProtonMail/go-cryptos).

The current implementation is using the `ProtonMail/go-crypto`-fork, but that
is an implementation-detail, and not exposed in any interfaces, so it might change
in the future.

## Local development

Run `docker compose run --rm --service-ports golang-devtools validate` to run
linting, vetting and tests.

## Releasing

This project uses `goreleaser` for releasing.
To release new version, follow the following steps:

1. Install [goreleaser](https://goreleaser.com/install/) locally.
2. Create a [GitHub Personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-personal-access-token-classic)
  with the scope `public_repo`, and authorize it to access the `coopnorge`
  GitHub organization with SSO. This token can be re-used during the lifetime of
  the token.
3. [Create a GPG signing-key](https://docs.github.com/en/authentication/managing-commit-signature-verification/generating-a-new-gpg-key)
  using the RSA or DSA algorithm (only these algorithms are supported due to
  a limitation in the [Terraform Public Provider Registry](https://developer.hashicorp.com/terraform/registry/providers/publishing#publishing-to-the-registry)).
  This GPG key can be re-used during the lifetime of the key.
4. A GitHub Organization Admin at Coop must upload the public part of the GPG-key
  from step 3 to the [Terraform Provider Registry](https://registry.terraform.io/settings/gpg-keys)
  under the `coopnorge` namespace. This must be done for every new key.
5. Create a new signed tag with git:

  ```shell
  git tag -a v0.2.0 --sign -m "Release v0.2.0"
  ```

6. Push the tag to GitHub

  ```shell
  git push origin v0.2.0
  ```

7. Run `goreleaser` to create a GitHub Release (substitute with your proper GitHub
  token and GPG fingerprint):

  ```shell
  GITHUB_TOKEN=githubtoken GPG_FINGERPRINT=gpgfingerprint goreleaser release --clean
  ```

8. Go to the newly created [GitHub release](https://github.com/coopnorge/terraform-provider-opengpg/releases/tag/v0.2.0),
  to verify that the changelog and artifacts looks correct.

9. Wait for the webhooks to kick in, and the new release should be available at
  [Terraform Provider Registry](https://registry.terraform.io/providers/coopnorge/opengpg)

## Configuring your repository

You can change settings in the repository by manupilating configuration file
in `.pallet/gitconfig.yaml`

The full definition of the api can be found [here][gitconfig-api-ref]

CI of the changes in the repositoryconfig is not configured yet. For now you
can use the [argocd app status][argocd-app-ref] to find the status and
potential errors.

[gitconfig-api-ref]: https://github.com/coopnorge/cloud-platform-apis/blob/main/cloud-platform-apis/templates/repositoryconfig.github.coop.no/definition.yaml
[argocd-app-ref]:  https://argocd.internal.coop/applications?search=pallet-terraform-provider-opengpg&showFavorites=false&proj=&sync=&autoSync=&health=&namespace=&cluster=&labels=
