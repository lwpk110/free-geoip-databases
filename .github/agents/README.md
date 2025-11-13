# Custom Agents

This directory contains custom agent configurations for GitHub Copilot.

## Available Agents

### Senior Coding Engineer (`senior-coding-engineer.yml`)

A highly experienced senior development engineer specializing in coding tasks.

**Expertise:**
- Primary Languages: Go, Python, JavaScript/TypeScript, Shell/Bash
- Frameworks: GitHub Actions, Docker, Git, REST APIs
- Specializations: Software Architecture, Code Optimization, Testing, DevOps, CI/CD

**Best For:**
- Implementing new features
- Bug fixing and debugging
- Code refactoring and optimization
- Writing and improving tests
- API development
- Script automation
- Performance improvements
- Security enhancements

**Example Usage:**

When working with GitHub Copilot in this repository, you can leverage this custom agent for coding tasks by mentioning coding-related requests. The agent is specifically tuned for:

1. **Go Development**: Writing and maintaining Go code for query tools and tests
2. **Shell Scripting**: Creating and improving automation scripts
3. **CI/CD**: Enhancing GitHub Actions workflows
4. **Database Tools**: Working with MMDB format databases
5. **Documentation**: Improving code examples and documentation

**Project-Specific Context:**

This agent is configured with specific knowledge about this repository:
- Focus: Free GeoIP databases with automatic updates
- Technologies: Go, Shell scripts, GitHub Actions, MMDB databases
- Key areas: Query tools, download automation, testing, CI/CD

## How Custom Agents Work

Custom agents are specialized configurations that enhance GitHub Copilot's capabilities for specific tasks or domains. They provide:

1. **Domain Expertise**: Specialized knowledge in specific areas
2. **Context Awareness**: Understanding of project-specific requirements
3. **Best Practices**: Adherence to coding standards and conventions
4. **Focused Assistance**: Targeted help for particular types of tasks

## Agent Configuration Structure

Each agent configuration includes:

- **name**: The agent's display name
- **description**: What the agent does
- **metadata**: Version, author, and categorization information
- **expertise**: Languages, tools, and specializations
- **behavior**: Coding style, problem-solving approach, and quality standards
- **task_preferences**: Types of tasks the agent excels at
- **communication**: How the agent communicates and documents
- **project_context**: Specific knowledge about this repository
- **instructions**: Detailed guidelines for the agent's operation
- **example_prompts**: Sample prompts that work well with the agent

## Creating New Agents

To create a new custom agent:

1. Create a new YAML file in this directory
2. Follow the structure shown in existing agents
3. Define the agent's expertise and capabilities
4. Customize behavior and communication style
5. Add project-specific context if needed
6. Test with example prompts

## Notes

- Custom agents enhance but don't replace GitHub Copilot's core capabilities
- Agents work best when given clear, specific tasks
- Multiple agents can be defined for different specializations
- Keep agent configurations up to date with project changes
