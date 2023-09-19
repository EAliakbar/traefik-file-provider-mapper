package internal

import (
	"fmt"
	"strings"

	"github.com/icza/dyno"
)

type MapperFunction = func(*Configuration)

func filterHttpRouters() MapperFunction {
	return func(c *Configuration) {
		HttpRouters, err := dyno.GetMapS(c.HTTP, "routers")
		if err != nil {
			return
		}
		for routerName, routerConfig := range HttpRouters {
			rawEntryPoints, err := dyno.GetSlice(routerConfig, "entryPoints")
			if err != nil {
				continue
			}
			for _, v := range rawEntryPoints {
				if str, ok := v.(string); !ok {
					if str != "http" {
						continue
					}
					err = dyno.Delete(c.HTTP, "routers", routerName)
					if err != nil {
						fmt.Printf("Could not remove routers.%s\n", routerName)
					}
					break
				}
			}
		}
	}
}

func filterWwwRouters() MapperFunction {
	return func(c *Configuration) {
		HttpRouters, err := dyno.GetMapS(c.HTTP, "routers")
		if err != nil {
			return
		}
		for routerName, routerConfig := range HttpRouters {
			rule, err := dyno.GetString(routerConfig, "rule")
			if err != nil {
				continue
			}
			if strings.Contains(rule, "Host(`www.") {
				err = dyno.Delete(c.HTTP, "routers", routerName)
				if err != nil {
					fmt.Printf("Could not remove routers.%s\n", routerName)
				}
			}
		}
	}
}

func renameEndpoint(oldName string, newName string) MapperFunction {
	return func(c *Configuration) {
		HttpRouters, err := dyno.GetMapS(c.HTTP, "routers")
		if err != nil {
			return
		}
		for routerName, routerConfig := range HttpRouters {
			rawEntryPoints, err := dyno.GetSlice(routerConfig, "entryPoints")
			if err != nil {
				continue
			}
			newEntryPoints := make([]interface{}, len(rawEntryPoints))
			for i, v := range rawEntryPoints {
				if str, ok := v.(string); ok {
					if str == oldName {
						newEntryPoints[i] = newName
					} else {
						newEntryPoints[i] = str
					}
				} else {
					// Handle the case where the element is not a string.
					fmt.Printf("Element at index %d is not a string: %v\n", i, v)
					newEntryPoints[i] = "" // Set a default value or handle the error as needed.
				}
			}
			err = dyno.Set(c.HTTP, newEntryPoints, "routers", routerName, "entryPoints")
			if err != nil {
				fmt.Printf("Could not set http.routers.%s.entryPoints\n", routerName)
			}
		}
	}
}

func renameTlsResolver(resolverName string) MapperFunction {
	return func(c *Configuration) {
		HttpRouters, err := dyno.GetMapS(c.HTTP, "routers")
		if err != nil {
			return
		}
		for routerName, routerConfig := range HttpRouters {
			_, err := dyno.GetString(routerConfig, "tls", "certresolver")
			if err != nil {
				continue
			}
			err = dyno.Set(c.HTTP, resolverName, "routers", routerName, "tls", "certresolver")
			if err != nil {
				fmt.Printf("Could not set http.routers.%s.tls.certresolver\n", routerName)
				return
			}
		}
	}
}

func filterCoolifyRouters() MapperFunction {
	return func(c *Configuration) {
		HttpRouters, err := dyno.GetMapS(c.HTTP, "routers")
		if err != nil {
			return
		}
		for routerName, _ := range HttpRouters {
			if strings.HasPrefix(routerName, "coolify") {
				fmt.Printf("Going to filter Route: %s\n", routerName)
				err := dyno.Delete(c.HTTP, routerName, "routers")
				if err != nil {
					fmt.Printf("Could Delete http.routers.%s\n", routerName)
				}
			}
		}
		return
	}
}

func filterCoolifyService() MapperFunction {
	return func(c *Configuration) {
		HttpServices, err := dyno.GetMapS(c.HTTP, "services")
		if err != nil {

			return
		}
		for serviceName, _ := range HttpServices {
			if strings.HasPrefix(serviceName, "coolify") {
				fmt.Printf("Going to filter Service: %s\n", serviceName)
				err := dyno.Delete(c.HTTP, serviceName, "services")
				if err != nil {
					fmt.Printf("Could Delete http.services.%s\n", serviceName)
				}
			}
		}
		return
	}
}

func (c *Configuration) Mapper(mc MapperConfig) {
	var actions []MapperFunction

	// filters
	if mc.removeHttpRouters {
		actions = append(actions, filterHttpRouters())
	}
	if mc.removeWwwRouters {
		actions = append(actions, filterWwwRouters())
	}

	// mappers
	if mc.newHttpEntrypointName != "" {
		actions = append(actions, renameEndpoint("http", mc.newHttpsEntrypointName))
	}
	if mc.newHttpsEntrypointName != "" {
		actions = append(actions, renameEndpoint("https", mc.newHttpsEntrypointName))
	}
	if mc.newCertResolver != "" {
		actions = append(actions, renameTlsResolver(mc.newCertResolver))
	}
	if mc.removeCoolify {
		actions = append(actions, filterCoolifyRouters())
		actions = append(actions, filterCoolifyService())
	}
	if mc.removeWwwMiddlewares {
		fmt.Println("Removing www Middlewares is not implemented")
	}
	if mc.ignoreMiddlewareSites != nil {
		fmt.Println("Ignoring Middlewares is not implemented")
	}

	for i, act := range actions {
		fmt.Printf("Action %d: \n", i)
		act(c)
	}
}
