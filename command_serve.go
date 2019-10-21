package main

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ory/hydra/sdk/go/hydra/client"
	"github.com/ory/hydra/sdk/go/hydra/client/admin"
	"github.com/ory/hydra/sdk/go/hydra/models"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

var hydraClient *client.OryHydra
var loginPage []byte
var loginFailedPage []byte

func loginRequest(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("login_challenge")

	contextLogger := log.WithFields(log.Fields{
		"method":    r.Method,
		"route":     "loginRequest",
		"challenge": challenge,
	})

	loginRequest, err := hydraClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge((challenge)))
	if err != nil {
		contextLogger.WithError(err).Error("GetLoginRequest failed")
		w.WriteHeader(500)
		return
	}

	contextLogger = contextLogger.WithField("skip", loginRequest.GetPayload().Skip)

	if loginRequest.GetPayload().Skip {
		acceptLoginRequest, err := hydraClient.Admin.AcceptLoginRequest(
			admin.NewAcceptLoginRequestParams().WithLoginChallenge(challenge).WithBody(&models.HandledLoginRequest{
				Subject: &loginRequest.GetPayload().Subject,
			}),
		)

		if err != nil {
			contextLogger.WithError(err).Error("AcceptLoginRequest failed")
			w.WriteHeader(500)
			return
		}

		contextLogger.WithField("subject", loginRequest.GetPayload().Subject).Info("Login accepted")
		http.Redirect(w, r, acceptLoginRequest.GetPayload().RedirectTo, 302)
		return
	}

	if r.Method != "POST" {
		contextLogger.Info("Login started")
		w.Write(loginPage)
		return
	}

	err = r.ParseForm()
	if err != nil {
		contextLogger.WithError(err).Error("ParseForm failed")
		w.WriteHeader(400)
		return
	}

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	contextLogger = contextLogger.WithField("subject", username)

	err = accountRepo.Exists(username, password)
	if err != nil {
		contextLogger.WithError(err).Info("Login failed")
		w.Write(loginFailedPage)
		return
	}

	acceptLoginRequest, err := hydraClient.Admin.AcceptLoginRequest(
		admin.NewAcceptLoginRequestParams().WithLoginChallenge(challenge).WithBody(&models.HandledLoginRequest{
			Subject: &username,
		}),
	)
	if err != nil {
		contextLogger.WithError(err).Error("AcceptLoginRequest failed")
		w.WriteHeader(500)
		return
	}

	contextLogger.Info("Login accepted")
	http.Redirect(w, r, acceptLoginRequest.GetPayload().RedirectTo, 302)
}

func consentRequest(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("consent_challenge")

	contextLogger := log.WithFields(log.Fields{
		"method":    r.Method,
		"route":     "consentRequest",
		"challenge": challenge,
	})

	contextLogger.Info("Consent started")

	consentRequest, err := hydraClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(challenge))
	if err != nil {
		contextLogger.WithError(err).Error("GetConsentRequest failed")
		w.WriteHeader(500)
		return
	}

	contextLogger = contextLogger.WithFields(log.Fields{
		"subject": consentRequest.GetPayload().Subject,
		"skip":    consentRequest.GetPayload().Skip,
	})

	acceptConsentRequest, err := hydraClient.Admin.AcceptConsentRequest(
		admin.NewAcceptConsentRequestParams().WithConsentChallenge(challenge).WithBody(&models.HandledConsentRequest{
			GrantedAudience: consentRequest.Payload.RequestedAudience,
			GrantedScope:    consentRequest.Payload.RequestedScope,
		}),
	)
	if err != nil {
		contextLogger.WithError(err).Error("AcceptConsentRequest failed")
		w.WriteHeader(500)
		return
	}

	contextLogger.Info("Consent accepted")
	http.Redirect(w, r, acceptConsentRequest.GetPayload().RedirectTo, 302)
}

func serveHandler() {
	var err error
	loginPage, err = ioutil.ReadFile(pathToLoginPage)
	if err != nil {
		log.WithField("page", pathToLoginPage).WithError(err).Fatal("Failed loading login page")
	}

	if pathToLoginFailedPage == "" {
		loginFailedPage = loginPage
	} else {
		loginFailedPage, err = ioutil.ReadFile(pathToLoginFailedPage)
		if err != nil {
			log.WithField("page", pathToLoginFailedPage).WithError(err).Fatal("Failed loading login failed page")
		}
	}

	adminURL, err := url.Parse(hydraURL)
	if err != nil {
		log.WithField("url", adminURL).WithError(err).Fatal("Failed parsing Hydra URL")
	}

	hydraClient = client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{adminURL.Scheme}, Host: adminURL.Host, BasePath: adminURL.Path})

	router := makeRouter([]routeDef{
		routeDef{"GET", "/login", "StartLoginRequest", loginRequest},
		routeDef{"POST", "/login", "FinishLoginRequest", loginRequest},
		routeDef{"GET", "/consent", "ConsentRequest", consentRequest},
	})

	if pathToStaticResources != "" {
		router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(pathToStaticResources))))
	}

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.UseHandler(router)

	log.WithField("port", port).Info("Listening")
	log.Fatal(http.ListenAndServe(":"+port, n))
}
