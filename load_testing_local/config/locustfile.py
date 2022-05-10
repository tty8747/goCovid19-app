import time
from locust import HttpUser, task, between


class Quickstart(HttpUser):
    wait_time = between(1, 5)

    @task
    def gocovid_test(self):
        self.client.request_name = "gocovid_test"
    #   self.client.get("http://web:4000/static/")
        self.client.get("http://web:4000/")
