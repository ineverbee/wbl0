package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(10, 30)

func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

type errorHandler func(http.ResponseWriter, *http.Request) error

func (f errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := f(w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			log.Printf("HTTP - %s", e)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

var base string = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">

	<link rel="icon" href="https://cdn.betterttv.net/emote/6038a2037c74605395f34895/3x">
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.3.1/dist/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
</head>
<body>
%s
</body>
</html>`

// GetDataPageHandler..
func GetDataPageHandler() errorHandler {
	return func(rw http.ResponseWriter, r *http.Request) error {
		param := mux.Vars(r)["id"]
		id, err := strconv.Atoi(param)
		if err != nil {
			return &StatusError{http.StatusBadRequest, fmt.Errorf("error: id is NaN")}
		}
		model, err := app.cache.Get(id)
		if err != nil {
			return &StatusError{http.StatusBadRequest, err}
		}
		tmpl := template.Must(template.New("data").Parse(fmt.Sprintf(base, `
		<h1>Data</h1>
		<table class="table">
			<thead>
				<tr>
				<th scope="col">#</th>
				<th scope="col">Key</th>
				<th scope="col">Value</th>
				</tr>
			</thead>
			<tbody>
				<tr>
					<th scope="row">1</th>
					<td>order_uid</td>
					<td>{{ .Order_uid}}</td>
				</tr>
				<tr>
					<th scope="row">2</th>
					<td>track_number</td>
					<td>{{ .Track_number}}</td>
				</tr>
				<tr>
					<th scope="row">3</th>
					<td>entry</td>
					<td>{{ .Entry}}</td>
				</tr>
				<tr>
					<th scope="row">4</th>
					<td>delivery</td>
					<td class="row">
						{{with .Delivery}}
						<table class="table table-sm table-hover">
							<thead class="bg-primary">
								<tr>
									<td scope="col">name</td>
									<td scope="col">phone</td>
									<td scope="col">zip</td>
									<td scope="col">city</td>
									<td scope="col">address</td>
									<td scope="col">region</td>
									<td scope="col">email</td>
								</tr>
							</thead>
							<tbody class="bg-secondary">
								<tr>
									<td>{{ .Name}}</td>
									<td>{{ .Phone}}</td>
									<td>{{ .Zip}}</td>
									<td>{{ .City}}</td>
									<td>{{ .Address}}</td>
									<td>{{ .Region}}</td>
									<td>{{ .Email}}</td>
								</tr>
							</tbody>
						</table>
						{{end}}							
					</td>
				</tr>
				<tr>
					<th scope="row">5</th>
					<td>payment</td>
					<td class="row">
						{{with .Payment}}
						<table class="table table-sm table-hover">
							<thead class="bg-primary">
								<tr>
									<td scope="col">transaction</td>
									<td scope="col">request_id</td>
									<td scope="col">currency</td>
									<td scope="col">provider</td>
									<td scope="col">bank</td>
									<td scope="col">amount</td>
									<td scope="col">payment_dt</td>
									<td scope="col">delivery_cost</td>
									<td scope="col">goods_total</td>
									<td scope="col">custom_fee</td>
								</tr>
							</thead>
							<tbody class="bg-secondary">
								<tr>
									<td>{{ .Transaction}}</td>
									<td>{{ .Request_id}}</td>
									<td>{{ .Currency}}</td>
									<td>{{ .Provider}}</td>
									<td>{{ .Bank}}</td>
									<td>{{ .Amount}}</td>
									<td>{{ .Payment_dt}}</td>
									<td>{{ .Delivery_cost}}</td>
									<td>{{ .Goods_total}}</td>
									<td>{{ .Custom_fee}}</td>
								</tr>
							</tbody>
						</table>
						{{end}}
					</td>
				</tr>
				<tr>
					<th scope="row">6</th>
					<td>items</td>
					<td class="row">{{with .Items}}
							<table class="table table-sm table-hover">
								<thead class="bg-primary">
									<tr>
										<td scope="col">track_number</td>
										<td scope="col">rid</td>
										<td scope="col">name</td>
										<td scope="col">size</td>
										<td scope="col">brand</td>
										<td scope="col">chrt_id</td>
										<td scope="col">price</td>
										<td scope="col">sale</td>
										<td scope="col">total_price</td>
										<td scope="col">nm_id</td>
										<td scope="col">status</td>
									</tr>
								</thead>
								<tbody class="bg-secondary">
							{{range .}}
									<tr>
										<td>{{ .Track_number}}</td>
										<td>{{ .Rid}}</td>
										<td>{{ .Name}}</td>
										<td>{{ .Size}}</td>
										<td>{{ .Brand}}</td>
										<td>{{ .Chrt_id}}</td>
										<td>{{ .Price}}</td>
										<td>{{ .Sale}}</td>
										<td>{{ .Total_price}}</td>
										<td>{{ .Nm_id}}</td>
										<td>{{ .Status}}</td>
									</tr>
							{{end}}
								</tbody>
							</table>
						{{end}}
					</td>
				</tr>
				<tr>
					<th scope="row">7</th>
					<td>locale</td>
					<td>{{ .Locale}}</td>
				</tr>
				<tr>
					<th scope="row">8</th>
					<td>internal_signature</td>
					<td>{{ .Internal_signature}}</td>
				</tr>
				<tr>
					<th scope="row">9</th>
					<td>customer_id</td>
					<td>{{ .Customer_id}}</td>
				</tr>
				<tr>
					<th scope="row">10</th>
					<td>delivery_service</td>
					<td>{{ .Delivery_service}}</td>
				</tr>
				<tr>
					<th scope="row">11</th>
					<td>shardkey</td>
					<td>{{ .Shardkey}}</td>
				</tr>
				<tr>
					<th scope="row">12</th>
					<td>oof_shard</td>
					<td>{{ .Oof_shard}}</td>
				</tr>
				<tr>
					<th scope="row">13</th>
					<td>sm_id</td>
					<td>{{ .Sm_id}}</td>
				</tr>
				<tr>
					<th scope="row">14</th>
					<td>date_created</td>
					<td>{{ .Date_created}}</td>
				</tr>
			</tbody>
		</table>`)))

		tmpl.Execute(rw, model)
		return nil
	}
}

// GetHomePageHandler..
func GetHomePageHandler() errorHandler {
	return func(rw http.ResponseWriter, r *http.Request) error {
		tmpl := template.Must(template.New("home").Parse(fmt.Sprintf(base, `
		<h1>Enter ID</h1>
		<form method="POST">
			<label>ID:</label><br />
			<input type="number" name="id"><br />
			<input type="submit">
		</form>`)))

		if r.Method != http.MethodPost {
			tmpl.Execute(rw, nil)
			return nil
		}

		_, err := strconv.Atoi(r.FormValue("id"))

		if err != nil {
			return &StatusError{http.StatusBadRequest, fmt.Errorf("error: id is NaN")}
		}

		http.Redirect(rw, r, "/data/"+r.FormValue("id"), http.StatusFound)
		return nil
	}
}
