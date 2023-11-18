function search() {
  // Declare variables 
  const input = document.getElementById("beerSearch");
  const filter = input.value.toUpperCase();
  const table = document.getElementById("beerTable");
  tr = table.getElementsByTagName("tr");

  // Loop through all table rows, and hide those who don't match the search query
  let brand
  let name
  for (i = 0; i < tr.length; i++) {
    brand = tr[i].getElementsByTagName("td")[0];
    name = tr[i].getElementsByTagName("td")[1];
    if (brand && name) {
      brandValue = brand.textContent || brand.innerText;
      nameValue = name.textContent || name.innerText;
      if (brandValue.toUpperCase().indexOf(filter) > -1 || nameValue.toUpperCase().indexOf(filter) > -1) {
        tr[i].style.display = "";
      } else {
        tr[i].style.display = "none";
      }
    } 
  }
}
