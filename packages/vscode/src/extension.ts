import * as vscode from 'vscode';
import fetch from 'node-fetch';

let statusBarItem: vscode.StatusBarItem;

export function activate(context: vscode.ExtensionContext) {
	console.log('GitHub Browser extension is now active');

	// 创建状态栏项
	statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Right, 100);
	statusBarItem.text = "$(github) GitHub Browser";
	statusBarItem.tooltip = "Open GitHub repository or PR";
	statusBarItem.command = 'github-browser.openRepository';
	context.subscriptions.push(statusBarItem);
	statusBarItem.show();

	// 检查服务状态
	checkServiceStatus();

	// 注册命令
	context.subscriptions.push(
		vscode.commands.registerCommand('github-browser.openRepository', openRepository)
	);

	context.subscriptions.push(
		vscode.commands.registerCommand('github-browser.openFromClipboard', openFromClipboard)
	);

	context.subscriptions.push(
		vscode.commands.registerCommand('github-browser.openPR', openPullRequest)
	);

	context.subscriptions.push(
		vscode.commands.registerCommand('github-browser.openConfig', openConfig)
	);
}

async function checkServiceStatus() {
	const config = vscode.workspace.getConfiguration('github-browser');
	const serviceUrl = config.get<string>('serviceUrl', 'http://localhost:9527');
	const autoCheck = config.get<boolean>('autoCheckService', true);

	if (!autoCheck) {
		return;
	}

	try {
		const response = await fetch(`${serviceUrl}/health`, { timeout: 2000 } as any);
		if (response.ok) {
			statusBarItem.text = "$(github) GitHub Browser";
			statusBarItem.backgroundColor = undefined;
			statusBarItem.tooltip = "GitHub Browser service is running";
		}
	} catch (error) {
		statusBarItem.text = "$(alert) GitHub Browser";
		statusBarItem.backgroundColor = new vscode.ThemeColor('statusBarItem.warningBackground');
		statusBarItem.tooltip = "GitHub Browser service is not running. Click to open configuration.";
		statusBarItem.command = 'github-browser.openConfig';
	}
}

async function openRepository() {
	const url = await vscode.window.showInputBox({
		prompt: 'Enter GitHub repository or PR URL',
		placeHolder: 'https://github.com/microsoft/vscode or https://github.com/microsoft/vscode/pull/123',
		validateInput: (value) => {
			if (!value) {
				return 'URL is required';
			}
			if (!value.includes('github.com')) {
				return 'Must be a GitHub URL';
			}
			return null;
		}
	});

	if (url) {
		await openGitHubURL(url);
	}
}

async function openFromClipboard() {
	const url = await vscode.env.clipboard.readText();

	if (!url) {
		vscode.window.showErrorMessage('Clipboard is empty');
		return;
	}

	if (!url.includes('github.com')) {
		vscode.window.showErrorMessage('Clipboard does not contain a GitHub URL');
		return;
	}

	await openGitHubURL(url);
}

async function openPullRequest() {
	// 输入 owner/repo
	const repo = await vscode.window.showInputBox({
		prompt: 'Enter repository (owner/repo)',
		placeHolder: 'microsoft/vscode'
	});

	if (!repo) {
		return;
	}

	// 输入 PR 号
	const prNumber = await vscode.window.showInputBox({
		prompt: 'Enter PR number',
		placeHolder: '12345',
		validateInput: (value) => {
			if (!value) {
				return 'PR number is required';
			}
			if (!/^\d+$/.test(value)) {
				return 'Must be a number';
			}
			return null;
		}
	});

	if (!prNumber) {
		return;
	}

	const url = `https://github.com/${repo}/pull/${prNumber}`;
	await openGitHubURL(url);
}

async function openGitHubURL(url: string) {
	const config = vscode.workspace.getConfiguration('github-browser');
	const serviceUrl = config.get<string>('serviceUrl', 'http://localhost:9527');

	// 显示进度
	await vscode.window.withProgress({
		location: vscode.ProgressLocation.Notification,
		title: "Opening GitHub repository...",
		cancellable: false
	}, async (progress) => {
		try {
			progress.report({ message: "Contacting service..." });

			const response = await fetch(`${serviceUrl}/open`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					url: url,
					ide: 'code'
				}),
				timeout: 60000 // 60 秒超时
			} as any);

			if (!response.ok) {
				const error = await response.json();
				throw new Error(error.message || 'Failed to open repository');
			}

			const result: any = await response.json();

			progress.report({ message: "Opening in VS Code..." });

			// 打开仓库
			if (result.path) {
				const uri = vscode.Uri.file(result.path);
				await vscode.commands.executeCommand('vscode.openFolder', uri, false);
			}

			vscode.window.showInformationMessage('✅ Repository opened successfully!');

		} catch (error: any) {
			console.error('Error opening repository:', error);

			if (error.code === 'ECONNREFUSED' || error.type === 'request-timeout') {
				vscode.window.showErrorMessage(
					'GitHub Browser service is not running. Please start the service first.',
					'Open Configuration'
				).then(selection => {
					if (selection === 'Open Configuration') {
						openConfig();
					}
				});
			} else {
				vscode.window.showErrorMessage(`Failed to open repository: ${error.message}`);
			}
		}
	});
}

async function openConfig() {
	const config = vscode.workspace.getConfiguration('github-browser');
	const serviceUrl = config.get<string>('serviceUrl', 'http://localhost:9527');

	const action = await vscode.window.showInformationMessage(
		`GitHub Browser Service URL: ${serviceUrl}`,
		'Open Settings',
		'Test Connection',
		'Open Service README'
	);

	if (action === 'Open Settings') {
		vscode.commands.executeCommand('workbench.action.openSettings', 'github-browser');
	} else if (action === 'Test Connection') {
		try {
			const response = await fetch(`${serviceUrl}/health`, { timeout: 2000 } as any);
			if (response.ok) {
				const data: any = await response.json();
				vscode.window.showInformationMessage(`✅ Service is running (version ${data.version})`);
				checkServiceStatus();
			}
		} catch (error) {
			vscode.window.showErrorMessage(
				'❌ Cannot connect to service. Make sure it is running:\n\n' +
				'1. Install: cd service && ./install.sh\n' +
				'2. Or run manually: ./github-browser-service'
			);
		}
	} else if (action === 'Open Service README') {
		vscode.env.openExternal(vscode.Uri.parse('https://github.com/your-repo/github-browser/tree/main/service'));
	}
}

export function deactivate() {
	if (statusBarItem) {
		statusBarItem.dispose();
	}
}
