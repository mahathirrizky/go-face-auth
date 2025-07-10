export function getSubdomain() {
    const host = window.location.hostname; // e.g., "owner.localhost" or "www.owner.localhost"

    const parts = host.split('.'); // Split the hostname by '.'
  
    // Check if there are at least 2 parts (domain and TLD)
    if (parts.length >= 2) {
        // If the first part is 'www', return the second part as the subdomain
        if (parts[0] === 'www') {
            return parts.length > 2 ? parts[1] : null; // Return the second part if there's a subdomain
        }
        return parts[0]; // Return the first part as the subdomain if it's not 'www'
    } else {
        return null; // Return null if no subdomain
    }
}

export function getBaseDomain() {
    const host = window.location.hostname;
    const parts = host.split('.');

    // Handle localhost specifically
    if (host === 'localhost') {
        return 'localhost';
    }

    // For domains like example.com, admin.example.com, example.co.id, admin.example.co.id
    // We want to get the last two parts for most common TLDs (e.g., .com, .org, .net)
    // or the last three parts for multi-part TLDs (e.g., .co.id, .com.au)
    // This is a simplified approach and might not cover all edge cases of TLDs.
    // A more robust solution might involve a list of known multi-part TLDs.
    if (parts.length > 2 && (parts[parts.length - 2] + '.' + parts[parts.length - 1]).includes('.')) {
        // Check for common multi-part TLDs (e.g., .co.id, .com.au)
        const multiPartTlds = ['co.id', 'com.au', 'org.uk', 'gov.uk']; // Add more as needed
        const lastTwoParts = parts[parts.length - 2] + '.' + parts[parts.length - 1];
        if (multiPartTlds.includes(lastTwoParts)) {
            return parts.slice(parts.length - 3).join('.'); // e.g., example.co.id
        }
    }
    
    // Default to the last two parts (e.g., example.com from admin.example.com)
    return parts.slice(parts.length - 2).join('.');
}