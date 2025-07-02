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