function myFunction(weight, id) {
  const body = {
    weight: parseFloat(weight.value),
  }
  fetch('/account/privilege/' +  id, {
    method: 'PUT',
    mode: 'cors',
    cache: 'no-cache',
    credentials: 'same-origin',
    body: JSON.stringify(body),
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  })
  .then(response => response.json())
  .then(data => {
    if (data?.error) {
      console.log(data.error);
      //message.error(data.error);
      return
    }
  })
  .catch((error) => {
    console.log(error);
    //message.error(error);
  });
}

function search() {
  // Declare variables 
  var input, filter, table, tr, td, i, txtValue;
  input = document.getElementById("myInput");
  filter = input.value.toUpperCase();
  table = document.getElementById("myTable");
  tr = table.getElementsByTagName("tr");

  // Loop through all table rows, and hide those who don't match the search query
  for (i = 0; i < tr.length; i++) {
    td = tr[i].getElementsByTagName("td")[0];
    if (td) {
      txtValue = td.textContent || td.innerText;
      if (txtValue.toUpperCase().indexOf(filter) > -1) {
        tr[i].style.display = "";
      } else {
        tr[i].style.display = "none";
      }
    } 
  }
}
