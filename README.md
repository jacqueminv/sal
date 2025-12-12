With `sal` you can fetch your own Strava activities from the past year straight from your comfy terminal.

Once its setup is complete, `sal` gladly fetches those activites with no further interaction.

# Pre-requirements

You'll need a [Strava API Application](https://www.strava.com/settings/api).

# Installation

`sal` binaries are available under the [release](https://github.com/jacqueminv/sal/releases) section. Simply download the one according to your platform of choice, unpack and that's it.

Alternatively, you are welcome to download `sal`'s code source and build it on your own machine. You'll need [golang](https://go.dev/dl/) for that.

# Launch `sal`

As simple as:

```sh
SAL_CLIENT_ID=<your_client_id> SAL_CLIENT_SECRET_ID=<your_client_secret_id> sal
```

`sal` needs a proper authentication token to interact with the Strava API, that's the reason why on first launch, it's automatically going to open your default browser to ask for your permission to fetch such a token.

> [!NOTE] 
> That requires obviously a desktop environment. Note though that once you have this token at disposal, you may simply transfer that file to another computer where you do not have a desktop environment and run `sal` without interactions.

If you grant that access, it's going to store that token locally and will be reused on further sessions.

In case you do not trust `sal` and deny access to your Strava account, it ends there.
