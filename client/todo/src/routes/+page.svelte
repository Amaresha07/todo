<script lang="ts">
	import { onMount } from 'svelte';

	let todos = [];
	let newTodoText = '';
	let errorMessage = '';
	let isLoading = false;
	let isEditing = false;
	let editingTodo: Todo | null = null;

	interface Todo {
		id: number;
		text: string;
		completed: boolean;
	}

	// Fetch todos from the API
	const fetchTodos = async () => {
		isLoading = true;
		try {
			const response = await fetch('http://localhost:8000/todos');
			if (!response.ok) throw new Error('Failed to fetch todos');
			todos = await response.json();
		} catch (error) {
			errorMessage = error.message;
		} finally {
			isLoading = false;
		}
	};

	// Add or Edit todo
	const saveTodo = async () => {
		if (!newTodoText.trim()) return; // Prevent adding empty todos
		const todo = { text: newTodoText, completed: false };

		let url = 'http://localhost:8000/todos/add';
		let method = 'POST';
		if (isEditing && editingTodo) {
			url = 'http://localhost:8000/todos/update';
			method = 'PUT';
			todo.id = editingTodo.id;
		}

		try {
			const response = await fetch(url, {
				method,
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(todo)
			});
			if (!response.ok) throw new Error('Failed to save todo');
			await fetchTodos(); // Refresh the list after adding or editing
			resetTodoForm(); // Reset form
		} catch (error) {
			errorMessage = error.message;
		}
	};

	// Toggle todo completion
	const toggleTodo = async (id: number) => {
		const todo = todos.find((t) => t.id === id);
		if (!todo) return;
		todo.completed = !todo.completed;

		try {
			const response = await fetch('http://localhost:8000/todos/update', {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(todo)
			});
			if (!response.ok) throw new Error('Failed to update todo');
			await fetchTodos(); // Refresh the list after toggling
		} catch (error) {
			errorMessage = error.message;
		}
	};

	// Delete a todo
	const deleteTodo = async (id: number) => {
		try {
			const response = await fetch('http://localhost:8000/todos/delete', {
				method: 'DELETE',

				body: JSON.stringify({ id })
			});
			if (!response.ok) throw new Error('Failed to delete todo');
			await fetchTodos(); // Refresh the list after deleting
		} catch (error) {
			errorMessage = error.message;
		}
	};

	// Start editing a todo
	const editTodo = (todo: Todo) => {
		isEditing = true;
		editingTodo = todo;
		newTodoText = todo.text;
	};

	// Reset the form
	const resetTodoForm = () => {
		isEditing = false;
		editingTodo = null;
		newTodoText = '';
	};

	onMount(fetchTodos); // Fetch todos on page load
</script>

<main class="flex flex-col items-center bg-gray-50 min-h-screen py-8">
	<h1 class="text-3xl font-semibold mb-6">Todo List</h1>

	<!-- Error Message -->
	{#if errorMessage}
		<div class="bg-red-200 text-red-600 p-4 mb-4 rounded-md">
			{errorMessage}
		</div>
	{/if}

	<!-- Loading Spinner -->
	{#if isLoading}
		<div class="flex justify-center items-center mb-6">
			<div
				class="spinner-border animate-spin inline-block w-12 h-12 border-4 border-blue-500 border-t-transparent rounded-full"
				role="status"
			>
				<span class="visually-hidden">Loading...</span>
			</div>
		</div>
	{/if}

	<!-- Todo Input Form -->
	<div class="w-full max-w-lg bg-white p-6 rounded-lg shadow-md">
		<div class="mb-4 flex items-center">
			<input
				id="todo-input"
				type="text"
				class="flex-grow p-3 border border-gray-300 rounded-l-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
				placeholder={isEditing ? 'Edit your todo...' : 'Add a new todo...'}
				bind:value={newTodoText}
			/>
			<button
				class="bg-blue-500 text-white p-3 rounded-r-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
				on:click={saveTodo}
			>
				{isEditing ? 'Save' : 'Add'}
			</button>
		</div>
	</div>

	<!-- Todo List -->
	{#if todos.length === 0 && !isLoading}
		<p class="text-gray-500 mb-4">No todos yet! Start by adding one.</p>
	{/if}

	<ul class="w-full max-w-lg space-y-4">
		{#each todos as todo}
			<li
				class="flex items-center justify-between bg-white p-4 rounded-lg shadow-md hover:shadow-lg transition-all duration-300"
			>
				<div class="flex items-center">
					<input
						type="checkbox"
						checked={todo.completed}
						class="w-5 h-5 text-blue-600 rounded focus:ring-blue-500"
						on:change={() => toggleTodo(todo.id)}
					/>
					<span class={`ml-3 ${todo.completed ? 'line-through text-gray-500' : 'text-gray-800'}`}>
						{todo.text}
					</span>
				</div>
				<div class="flex space-x-2">
					<button
						class="text-yellow-500 hover:text-yellow-700 focus:outline-none"
						on:click={() => editTodo(todo)}
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-6 w-6"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M5 13l4 4L19 7m2 1v6a2 2 0 002 2h-6m0-8V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3"
							/>
						</svg>
					</button>
					<button
						class="text-red-500 hover:text-red-700 focus:outline-none"
						on:click={() => deleteTodo(todo.id)}
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-6 w-6"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
							/>
						</svg>
					</button>
				</div>
			</li>
		{/each}
	</ul>
</main>

<style>
	.spinner-border {
		border-width: 4px;
		border-top-color: transparent;
	}
</style>
