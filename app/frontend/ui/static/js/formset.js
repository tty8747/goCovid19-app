let dateFrom = document.getElementById('dateFrom')
let dateTo = document.getElementById('dateTo')

let countrySelect = document.getElementById('countrySelect')
let radioDD = document.getElementsByName('radioDD')

document.addEventListener('submit', () => {
	sessionStorage.setItem('date_from', dateFrom.value)
	sessionStorage.setItem('date_to', dateTo.value)
//  if (countrySelect.value == "") {
//	  countrySelect.value = "default"
//  }
	sessionStorage.setItem('country_select', countrySelect.value)
	radioDD.forEach(button => {
		if (button.checked) {
			sessionStorage.setItem('radio_dd', button.value)
		}
	})
})

dateFrom.value = sessionStorage.getItem('date_from')
dateTo.value = sessionStorage.getItem('date_to')
if (sessionStorage.getItem('country_select') == null) {
  countrySelect.value = 'default'
} else {
  countrySelect.value = sessionStorage.getItem('country_select')
}

let radioValue = sessionStorage.getItem('radio_dd')
let radioButton = document.querySelector(`input[value = "${radioValue}"]`)
radioButton.checked = true
