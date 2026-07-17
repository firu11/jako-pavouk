export type OperatingSystem = 'macos' | 'windows' | 'linux' | 'ios' | 'android' | 'unknown';

export function getOperatingSystem(): OperatingSystem {
    const { userAgent, platform, maxTouchPoints } = window.navigator;

    if (/iPhone|iPad|iPod/i.test(userAgent) || (platform === 'MacIntel' && maxTouchPoints > 1)) return 'ios';
    if (/Android/i.test(userAgent)) return 'android';
    if (/macOS|Macintosh|MacIntel|MacPPC|Mac68K|darwin/i.test(userAgent)) return 'macos';
    if (/Windows|Win32|Win64|WinCE/i.test(userAgent)) return 'windows';
    if (/Linux/i.test(userAgent)) return 'linux';

    return 'unknown';
}
