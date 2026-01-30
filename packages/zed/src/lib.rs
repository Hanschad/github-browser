use serde::{Deserialize, Serialize};
use std::fs;
use zed_extension_api::{self as zed, Command, Result};

struct GitHubBrowserExtension {
    service_url: String,
}

#[derive(Serialize)]
struct OpenRequest {
    url: String,
    ide: String,
}

#[derive(Deserialize)]
struct OpenResponse {
    status: String,
    message: String,
    path: Option<String>,
}

impl zed::Extension for GitHubBrowserExtension {
    fn new() -> Self {
        Self {
            service_url: "http://localhost:9527".to_string(),
        }
    }

    fn language_server_command(
        &mut self,
        _language_server_id: &zed::LanguageServerId,
        _worktree: &zed::Worktree,
    ) -> Result<Command> {
        Ok(Command {
            command: "".to_string(),
            args: vec![],
            env: Default::default(),
        })
    }
}

impl GitHubBrowserExtension {
    fn open_github_url(&self, url: &str) -> Result<()> {
        let request = OpenRequest {
            url: url.to_string(),
            ide: "zed".to_string(),
        };

        let client = reqwest::blocking::Client::new();
        let response = client
            .post(format!("{}/open", self.service_url))
            .json(&request)
            .send()
            .map_err(|e| format!("Failed to send request: {}", e))?;

        if !response.status().is_success() {
            return Err(format!("Service returned error: {}", response.status()).into());
        }

        let result: OpenResponse = response
            .json()
            .map_err(|e| format!("Failed to parse response: {}", e))?;

        if result.status != "ok" {
            return Err(result.message.into());
        }

        Ok(())
    }
}

zed::register_extension!(GitHubBrowserExtension);
