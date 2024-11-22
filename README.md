# GPG Terraform provider

The GPG provider allows to generate GPG encrypted message in ASCII-armored
format using Terraform. It is a fork of
[invidian/terraform-provider-gpg](https://github.com/invidian/terraform-provider-gpg).
The provider, and it's resources has been renamed from `gpg` to `opengpg`, to
make a clearer distinction between the two providers.

The provider has been forked to add support for newer encryption-algorithms,
such as [Curve25519](https://en.wikipedia.org/wiki/Curve25519).

The original provider was using `golang.org/x/crypto/openpgp`, but that has
been marked frozen and deprecated, and it does not support Curve25519. The Go
team is promoting using community-forks such as
[keybase/go-crypto](https://github.com/keybase/go-crypto) or
[ProtonMail/go-crypto](https://github.com/ProtonMail/go-cryptos).

The current implementation is using the `ProtonMail/go-crypto`-fork, but that
is an implementation-detail, and not exposed in any interfaces, so it might
change in the future.

## Local development

Run `docker compose run --rm --service-ports golang-devtools validate` to run
linting, vetting and tests.

## Releasing

This project uses `goreleaser` for releasing. To release new version, follow
the following steps:

1. Create a new signed tag with git:

```shell
git tag -a v0.2.0 --sign -m "Release v0.2.0"
```

2. Push the tag to GitHub

```shell
git push origin v0.2.0
```

3. A GitHub Actions should start running. You can find and monitor the run by
   going to the
   [GitHub Actions page](https://github.com/coopnorge/terraform-provider-opengpg/actions/workflows/release.yaml).

4. Go to the newly created
   [GitHub release](https://github.com/coopnorge/terraform-provider-opengpg/releases/tag/v0.2.0),
   to verify that the changelog and artifacts looks correct.

5. Wait for the Webhook from GitHub to Terraform Provider Registry to kick in,
   and the new release should be available at
   [Terraform Provider Registry](https://registry.terraform.io/providers/coopnorge/opengpg)
   a. If the webhook does not kick in after a certain time, a manual resync can
   be triggered from the
   [Terraform Provider Registry Settings](https://registry.terraform.io/providers/coopnorge/opengpg/latest/settings/resync).

### GPG Key for signing the release

The GitHub Release created by `goreleaser` must be signed with a GPG key that
is known to the Terraform Provider Registry. The following steps must be
performed to rotate the key, or if we want to publish another provider.

This GPG key was created by running `gpg --full-generate-key` and selecting the
key type `DSA (sign only)`, and key size 3072 bits, and a passphrase.

The private key was then exported by running
`gpg --armor --export-secret-keys <key-id-from-previous-command>`.

Then both the private key and the passphrase was uploaded to GitHub Secrets,
with the keys `GPG_PRIVATE_KEY` and `GPG_PASSPHRASE`, respectively.

The public key was then exported by running
`gpg --armor --export <key-id-from-previous-command>`, and then uploaded (by a
GitHub Organization Admin at Coop) to the
[Terraform Provider Registry](https://registry.terraform.io/settings/gpg-keys),
under the `coopnorge namespace.

## Configuring your repository

You can change settings in the repository by manipulating configuration file in
`.pallet/gitconfig.yaml`

The full definition of the api can be found [here][gitconfig-api-ref]

CI of the changes in the repositoryconfig is not configured yet. For now you
can use the [argocd app status][argocd-app-ref] to find the status and
potential errors.

[gitconfig-api-ref]:
  https://github.com/coopnorge/cloud-platform-apis/blob/main/cloud-platform-apis/templates/repositoryconfig.github.coop.no/definition.yaml
[argocd-app-ref]:
  https://argocd.internal.coop/applications?search=pallet-terraform-provider-opengpg&showFavorites=false&proj=&sync=&autoSync=&health=&namespace=&cluster=&labels=

```

```
