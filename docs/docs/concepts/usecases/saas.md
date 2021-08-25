---
title: Saas Product with Authentication and Authorization
---

This is an example architecture for a typical saas product. 
To illustrate it a fictional organization and project is taken.

## Example Case

The Timing Company has a product called Time.
They have two environments, the development and the production environment.
In this case the time uses authentication and authorizations from ZITADEL.
This means that the users and also their authorizations will be managed within ZITADEL.

![Architecture](/img/concepts/usecase/saas.png)

## Organization

An organization is the top level in ZITADEL. 
In an organization projects and users are managed by the organization.
You need at least one organization for your own company in our case "The Timing Company".

For your costumers you have different possibilities:
1. Your customer already owns an organization ZITADEL
2. Your customer creates a new organization in ZITADEL by itself
3. You create an organization for your customer (If you like to verify the domain, the customer has to do it)

:::info
Subscriptions are organization based. This means, that each organization can choose her own tier based on the needed features.
:::

## Project

The idea of projects is to have a vessel for all components who are closely related to each other.

In our use case we would have two different projects, for each environment one. So lets call it "Time Dev" and "Time Prod".
These projects should be created in "The Timing Company" organization, because this is the owner of the project.

In the project you will configure all your roles and applications (clients and APIs).

### Project Settings

If it should only be allowed to login to the application of the project if a user has an authorization, this can be configured on the project. (Check roles on authentication)

### Project Grant

To give a customer permissions to a project, a project grant to the customer is need (Search by domain).
It is also possible to only allow the customer to use specific roles.
As soon as a project grant exists, the customer will see the project in the granted projects section of his organization and will be able to authorize his own users to the given proejct.

## Authorizations

To give a user permission to a project an authorization is need.
All organizations which own the project or got it granted are able to authorize users.
It is also possible to authorize users outside the own company if the exact login name of the user is known.

## Project Login

There are some different use cases how the login should behave and look like:

1. Restrict Organization

With the primary domain scope the organization will be restricted to the requested domain, this means only users of the requestd organization will be able to login.
The private labeling (branding) and the login policy of the requested organization will trigger automatically.

:::note
More about the [Scopes](../../apis/openidoauth/scopes)
:::

2. Show Private Labeling (Branding) of the project organization

On the project can be configured, what kind of branding should be shown.
In the default the design of ZITADEL will be shown, but as soon as the user is identified, the policy of the users organization will be triggered.
If the setting is set to Ensure Project Resource Owner Setting, the private labeling of the project organization will always be triggered.
The last possibility is to show the private labeling of the project organization and as soon as the user is identitfied the user organization settings will be triggered.
For this the Allow User Resource Owner Setting should be set.
:::note
More about [Private Labeling](../../guides/customization/branding)
:::