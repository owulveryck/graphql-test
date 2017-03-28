import React, { Component } from 'react';
import './App.css';
import ReactTable from 'react-table'
import 'react-table/react-table.css'
import { ApolloClient, ApolloProvider, createNetworkInterface } from 'react-apollo';
import { gql, graphql } from 'react-apollo';

const columns = [{
  header: 'SKU',
  accessor: 'sku' // String-based value accessors!
}, {
  header: 'Location',
  accessor: 'location',
  sortable: true,
  width: 400,
  style: {
    textAlign: 'left'
  }
}, {
  header: 'Instance Type',
  accessor: 'instanceType'
}, {
  header: 'Operating System',
  accessor: 'operatingSystem'
}]

function ProductList({ loading, products }) {
  if (loading) {
    return <div>Loading</div>;
  } else {
    return (
      <div className="App">
      <ReactTable className="-striped -highlight"
      pivotBy={['location', 'instanceType']}
      data={products}
      columns={columns}
      />
      </div>
    );
  }
}

const allProducts = gql`
query products {
  products{
    sku
    location
    instanceType
    operatingSystem
  }
}
`

const ProductListWithData = graphql(allProducts, {
  props: ({data: { loading, products }}) => ({
    loading,
      products,
  }),
})(ProductList);

class App extends Component {
  constructor(props) {
    super(props);
    const networkInterface = createNetworkInterface({
      uri: 'http://localhost:8080/graphql'
    })

    this.client = new ApolloClient({
      networkInterface: networkInterface
    });
  }

  render() {
    return (
      <ApolloProvider client={this.client}>
        <ProductListWithData />
      </ApolloProvider>
    );
  }
}

export default App;
