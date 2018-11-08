package main

import (
	"github.com/go-gremlin/gremlin"
	"fmt"
)

func main()  {

	if err := gremlin.NewCluster("ws://192.168.1.66:8182/gremlin"); err != nil {
		// handle error
		fmt.Println(err)
	}
	//:remote connect tinkerpop.server conf/remote.yaml

	//data, err := gremlin.Query(`g.V().has("name", userName).valueMap()`).Bindings(gremlin.Bind{"userName": "john"}).Exec()
	//gremlin.Query(`graph.addVertex("name", "aaaaaaaaaaaaa")`).Exec()
	//data, _ := gremlin.Query(`g.V()`).Exec()
	//fmt.Println(data)
	//res,_ := gremlin.Query(`g.V().values('name')`).Exec()
	//res,_ := gremlin.Query(`graph.openManagement()`).Exec()
	//res,_ := gremlin.Query(`graph.openManagement().makeEdgeLabel('follow').multiplicity(MULTI).make()`).Exec()
	//res,_ := gremlin.Query(`graph.openManagement().makeEdgeLabel('follow').multiplicity(MULTI).make()`).Exec()
	//data, err := gremlin.Query(`g.V().has('resumeID').limit(1).as('a').out('employment').in('employment').has('birth', gt(0)).dedup().as('b').where('a',neq('b')).where('a',eq('b')).by('birth').valueMap()`).Exec()

	//fmt.Println(res)



	//str := `graph.openManagement().makeEdgeLabel('follow').multiplicity(MULTI).make()`
	//str := `graph.openManagement().makeEdgeLabel('mother').multiplicity(MANY2ONE).make()`
	//str := `graph.openManagement().commit()`
	//str := ``
	//str := `g.V(1).as('a').out('created').in('created').where(neq('a')).addE('co-developer').from('a').property('year',2009)`
	//str := `g.V(3,4,5).aggregate('x').has('name','josh').as('a').select('x').unfold().hasLabel('software').addE('createdBy').to('a')`
	//str := `g.V().as('a').out('created').addE('createdBy').to('a').property('acl','public')`
	//str := `g.V(1).as('a').out('knows').addE('livesNear').from('a').property('year',2009).inV().inE('livesNear').values('year')`
	//str := `g.V().match(
    //             __.as('a').out('knows').as('b'),
    //             __.as('a').out('created').as('c'),
    //             __.as('b').out('created').as('c')).
    //           addE('friendlyCollaborator').from('a').to('b').
    //             property(id,23).property('project',select('c').values('name'))`

	//str :=  `g.E(23).valueMap()`
	//str := `g.V(has('name','marko').next()).addE('knows').to(has('name','peter').next())`
	//str := `g.V()`

	//s := `graph = TinkerFactory.createModern()`
	s := `g = graph.traversal()`
	res,err := gremlin.Query(s).Exec()
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(string(res))





	//建模型





}
