const tbody = document.getElementById('account-list-table-body');

loadAccounts();


async function loadAccounts() {
    try {
        const response = await fetch('http://localhost:8080/api/v1/accounts');
        const accounts = await response.json();

        tbody.innerHTML = '';
        accounts.forEach(account => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${account.id}</td>
                <td>${account.name}</td>
                <td>${account.currencyId}</td>
                <td>${account.balance}</td>
            `;
            tbody.appendChild(row);
        });
    } catch (error) {
        console.error('Error:', error);
    }   
}

