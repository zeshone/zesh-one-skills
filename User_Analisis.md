**tengo en mente que hicimos lo siguiente**
1- corregimos las skills de net8-apirest/* - OK
2- agregamos /testing-unit - OK
3- agregamos pruebas de formato para las skills con /tools - OK
4- iniciamos repositorio con Git, readme, gitignore, credits - OK

**que falta?**
1- separar skills por dominio
	- BackEnd: Todas las skills de net8-apirest/*
	- FrontEnd: las skills - nextjs-15, react-19, tailwind-4, typescript, zod-4, zustand-5
	- Github-pr
2- dejamos fuera como auxiliar
	- agents-files
	- rules-to-skill
3- Creamos archivos agents.md por dominio basandonos en las plantillas de agnets-files
	- BackEnd: agents.md -> auto-invoke skills de net8-apirest con sus condiciones de invocacion, adicional la invocacion para la skill de Github-pr
	- FrontEnd: agents.md -> auto-invoke skills basicamente la plantilla de prowler nos hace todo el trabajo solo acopla los que si tenemos y quita los que no, igual no olvides el Github-pr
nota: Lo que no se puede definir en agents.md en este momento solo coloca placehorders para que se definan en cada proyecto
